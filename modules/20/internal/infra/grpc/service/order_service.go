package service

import (
	"context"

	"github.com/joaqu1m/goexpert-labs/modules/20/internal/infra/grpc/pb"
	"github.com/joaqu1m/goexpert-labs/modules/20/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	order, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         order.ID,
		Price:      float32(order.Price),
		Tax:        float32(order.Tax),
		FinalPrice: float32(order.FinalPrice),
	}, nil
}
