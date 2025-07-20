package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/configs"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/event/rabbitmq"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/graph"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/grpc/pb"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/grpc/service"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/internal/infra/web/webserver"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/pkg/events"
	"github.com/joaqu1m/goexpert-labs/projects/clean-architecture/wired"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open(configs.Cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.Cfg.DBUser, configs.Cfg.DBPassword, configs.Cfg.DBSource, configs.Cfg.DBPort, configs.Cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()
	defer rabbitMQChannel.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.RegisterHandler("OrderListed", &rabbitmq.OrderListedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.RegisterHandler("OrderCreated", &rabbitmq.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := wired.NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := wired.NewListOrdersUseCase(db, eventDispatcher)

	webServer := webserver.NewWebServer(configs.Cfg.WebServerPort)

	webListOrdersHandler := wired.NewWebListOrdersHandler(db, eventDispatcher)
	webServer.AddHandler("/orders/list", webListOrdersHandler.List)
	webCreateOrderHandler := wired.NewWebCreateOrderHandler(db, eventDispatcher)
	webServer.AddHandler("/orders/create", webCreateOrderHandler.Create)

	fmt.Println("Starting web server on port:", configs.Cfg.WebServerPort)
	go webServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port:", configs.Cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.Cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port:", configs.Cfg.GraphQLServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.Cfg.GraphQLServerPort), nil)
}

func getRabbitMQChannel() *amqp091.Channel {
	rabbitMQChannel, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.Cfg.RabbitMQUser, configs.Cfg.RabbitMQPassword, configs.Cfg.RabbitMQSource, configs.Cfg.RabbitMQPort))
	if err != nil {
		panic(err)
	}

	ch, err := rabbitMQChannel.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
