package auction

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/joaqu1m/goexpert-labs/projects/auction-concurrency/configuration/logger"
	"github.com/joaqu1m/goexpert-labs/projects/auction-concurrency/internal/entity/auction_entity"
	"github.com/joaqu1m/goexpert-labs/projects/auction-concurrency/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection         *mongo.Collection
	auctionInterval    time.Duration
	stopAuctionChecker chan bool
	auctionMutex       *sync.Mutex
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection:         database.Collection("auctions"),
		auctionInterval:    getAuctionInterval(),
		stopAuctionChecker: make(chan bool),
		auctionMutex:       &sync.Mutex{},
	}

	go repo.auctionCloser()

	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}
	return duration
}

func (ar *AuctionRepository) auctionCloser() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ar.closeExpiredAuctions()
		case <-ar.stopAuctionChecker:
			return
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions() {
	ar.auctionMutex.Lock()
	defer ar.auctionMutex.Unlock()

	ctx := context.Background()

	currentTime := time.Now()
	expiredTimestamp := currentTime.Add(-ar.auctionInterval).Unix()

	filter := bson.M{
		"status": auction_entity.Active,
		"timestamp": bson.M{
			"$lte": expiredTimestamp,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status": auction_entity.Completed,
		},
	}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to close expired auctions", err)
		return
	}

	if result.ModifiedCount > 0 {
		logger.Info("Closed expired auctions", zap.Int64("count", result.ModifiedCount))
	}
}

func (ar *AuctionRepository) StopAuctionChecker() {
	close(ar.stopAuctionChecker)
}
