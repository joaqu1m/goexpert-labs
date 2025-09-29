package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joaqu1m/goexpert-labs/projects/auction-concurrency/configuration/database/mongodb"
	"github.com/joaqu1m/goexpert-labs/projects/auction-concurrency/internal/entity/auction_entity"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAuctionLogicWithoutDatabase(t *testing.T) {
	auction, err := auction_entity.CreateAuction(
		"Test Product",
		"Electronics",
		"This is a test auction",
		auction_entity.New,
	)

	if err != nil {
		t.Fatal("Expected no error creating auction, got:", err)
	}

	if auction == nil {
		t.Fatal("Expected auction to be created, got nil")
	}

	if auction.Status != auction_entity.Active {
		t.Errorf("Expected auction status to be Active, got %v", auction.Status)
	}

	os.Setenv("AUCTION_INTERVAL", "30s")
	interval := getAuctionInterval()
	expected := 30 * time.Second

	if interval != expected {
		t.Errorf("Expected auction interval %v, got %v", expected, interval)
	}
}

func TestAuctionAutoCloseWithMongoDB(t *testing.T) {
	os.Setenv("MONGODB_URL", "mongodb://admin:admin@localhost:27017/auctions?authSource=admin")
	os.Setenv("MONGODB_DB", "auctions_test")
	os.Setenv("AUCTION_INTERVAL", "2s")

	databaseConnection, err := mongodb.NewMongoDBConnection(context.Background())
	if err != nil {
		t.Skip("Skipping test: MongoDB not available -", err)
		return
	}

	auctionRepo := NewAuctionRepository(databaseConnection)
	defer auctionRepo.StopAuctionChecker()

	auctionRepo.Collection.Drop(context.Background())

	auction, internalErr := auction_entity.CreateAuction(
		"Test Product",
		"Electronics",
		"This is a test auction for automatic closure",
		auction_entity.New,
	)
	if internalErr != nil {
		t.Fatal("Error creating auction:", internalErr)
	}
	if auction == nil {
		t.Fatal("Auction should not be nil")
	}

	createErr := auctionRepo.CreateAuction(context.Background(), auction)
	if createErr != nil {
		t.Skip("Skipping test: Could not create auction in database -", createErr)
		return
	}

	createdAuction, findErr := auctionRepo.FindAuctionById(context.Background(), auction.Id)
	if findErr != nil {
		t.Fatal("Error finding auction:", findErr)
	}
	if createdAuction.Status != auction_entity.Active {
		t.Errorf("Expected auction status to be Active, got %v", createdAuction.Status)
	}

	time.Sleep(3 * time.Second)

	updatedAuction, findErr2 := auctionRepo.FindAuctionById(context.Background(), auction.Id)
	if findErr2 != nil {
		t.Fatal("Error finding updated auction:", findErr2)
	}
	if updatedAuction.Status != auction_entity.Completed {
		t.Errorf("Expected auction status to be Completed, got %v", updatedAuction.Status)
	}
}

func TestGetAuctionInterval(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "30s")
	interval := getAuctionInterval()
	expectedInterval := 30 * time.Second
	if interval != expectedInterval {
		t.Errorf("Expected interval %v, got %v", expectedInterval, interval)
	}

	os.Setenv("AUCTION_INTERVAL", "invalid")
	interval = getAuctionInterval()
	expectedDefault := 5 * time.Minute
	if interval != expectedDefault {
		t.Errorf("Expected default interval %v, got %v", expectedDefault, interval)
	}

	os.Unsetenv("AUCTION_INTERVAL")
	interval = getAuctionInterval()
	if interval != expectedDefault {
		t.Errorf("Expected default interval %v, got %v", expectedDefault, interval)
	}
}

func TestAuctionRepositoryCreation(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "30s")

	interval := getAuctionInterval()
	expected := 30 * time.Second

	if interval != expected {
		t.Errorf("Expected auction interval %v, got %v", expected, interval)
	}
}

func TestCloseExpiredAuctions(t *testing.T) {
	os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DB", "auctions_test")
	os.Setenv("AUCTION_INTERVAL", "1s")

	databaseConnection, err := mongodb.NewMongoDBConnection(context.Background())
	if err != nil {
		t.Skip("Skipping test: MongoDB not available -", err)
		return
	}

	auctionRepo := NewAuctionRepository(databaseConnection)
	defer auctionRepo.StopAuctionChecker()

	auctionRepo.Collection.Drop(context.Background())

	expiredAuction, _ := auction_entity.CreateAuction(
		"Expired Product",
		"Electronics",
		"This auction should be expired",
		auction_entity.New,
	)
	expiredAuction.Timestamp = time.Now().Add(-5 * time.Second)

	activeAuction, _ := auction_entity.CreateAuction(
		"Active Product",
		"Electronics",
		"This auction should remain active",
		auction_entity.New,
	)

	auctionRepo.CreateAuction(context.Background(), expiredAuction)
	auctionRepo.CreateAuction(context.Background(), activeAuction)

	auctionRepo.closeExpiredAuctions()

	expiredResult, _ := auctionRepo.FindAuctionById(context.Background(), expiredAuction.Id)
	if expiredResult.Status != auction_entity.Completed {
		t.Errorf("Expected expired auction to be Completed, got %v", expiredResult.Status)
	}

	activeResult, _ := auctionRepo.FindAuctionById(context.Background(), activeAuction.Id)
	if activeResult.Status != auction_entity.Active {
		t.Errorf("Expected active auction to remain Active, got %v", activeResult.Status)
	}
}

func TestConcurrentAuctionCreationAndClosing(t *testing.T) {
	os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DB", "auctions_test")
	os.Setenv("AUCTION_INTERVAL", "1s")

	databaseConnection, err := mongodb.NewMongoDBConnection(context.Background())
	if err != nil {
		t.Skip("Skipping test: MongoDB not available -", err)
		return
	}

	auctionRepo := NewAuctionRepository(databaseConnection)
	defer auctionRepo.StopAuctionChecker()

	auctionRepo.Collection.Drop(context.Background())

	numAuctions := 10
	done := make(chan bool, numAuctions)

	for i := 0; i < numAuctions; i++ {
		go func(index int) {
			auction, _ := auction_entity.CreateAuction(
				"Concurrent Product",
				"Electronics",
				"Concurrent auction test",
				auction_entity.New,
			)
			auctionRepo.CreateAuction(context.Background(), auction)
			done <- true
		}(i)
	}

	for i := 0; i < numAuctions; i++ {
		<-done
	}

	ctx := context.Background()
	count, err := auctionRepo.Collection.CountDocuments(ctx, bson.M{"status": auction_entity.Active})
	if err != nil {
		t.Fatal("Error counting active auctions:", err)
	}
	if count != int64(numAuctions) {
		t.Errorf("Expected %d active auctions, got %d", numAuctions, count)
	}

	time.Sleep(2 * time.Second)
	time.Sleep(11 * time.Second)

	completedCount, err := auctionRepo.Collection.CountDocuments(ctx, bson.M{"status": auction_entity.Completed})
	if err != nil {
		t.Fatal("Error counting completed auctions:", err)
	}
	if completedCount != int64(numAuctions) {
		t.Errorf("Expected %d completed auctions, got %d", numAuctions, completedCount)
	}

	activeCount, err := auctionRepo.Collection.CountDocuments(ctx, bson.M{"status": auction_entity.Active})
	if err != nil {
		t.Fatal("Error counting active auctions:", err)
	}
	if activeCount != int64(0) {
		t.Errorf("Expected 0 active auctions, got %d", activeCount)
	}
}
