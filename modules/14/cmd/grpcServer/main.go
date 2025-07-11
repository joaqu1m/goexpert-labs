package main

import (
	"database/sql"
	"net"

	"github.com/joaqu1m/goexpert-labs/modules/14/internal/pb"
	"github.com/joaqu1m/goexpert-labs/modules/14/internal/repository"
	"github.com/joaqu1m/goexpert-labs/modules/14/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(*categoryRepository)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
