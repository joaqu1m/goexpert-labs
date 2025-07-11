package service

import (
	"context"
	"errors"
	"io"

	"github.com/joaqu1m/goexpert-labs/modules/14/internal/pb"
	"github.com/joaqu1m/goexpert-labs/modules/14/internal/repository"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		CategoryRepository: repo,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryRepository.CreateCategory(req.Name, &req.Description)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, req *pb.Blank) (*pb.CategoryList, error) {
	categories, err := s.CategoryRepository.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return &pb.CategoryList{
		Categories: categories,
	}, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, req *pb.CategoryGetRequest) (*pb.Category, error) {
	categories, err := s.CategoryRepository.GetAllCategories()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if category.Id == req.Id {
			return category, nil
		}
	}

	return nil, errors.New("category not found")
}

func (s *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}

		if err != nil {
			if err == context.Canceled {
				return nil
			}
			return err
		}

		category, err := s.CategoryRepository.CreateCategory(req.Name, &req.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, category)
	}
}

func (s *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			if err == context.Canceled {
				return nil
			}
			return err
		}

		category, err := s.CategoryRepository.CreateCategory(req.Name, &req.Description)
		if err != nil {
			return err
		}

		if err := stream.Send(category); err != nil {
			return err
		}
	}
}
