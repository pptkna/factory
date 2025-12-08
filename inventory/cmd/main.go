package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	partApiV1 "github.com/pptkna/rocket-factory/inventory/internal/api/part/v1"
	partRepository "github.com/pptkna/rocket-factory/inventory/internal/repository/part"
	partService "github.com/pptkna/rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051

	host     = "localhost"
	port     = "27017"
	user     = "inventory-service-user"
	password = "inventory-service-password"
	dbname   = "inventory-service"

	dbURI = "mongodb://inventory-service-user:inventory-service-password@localhost:27017/inventory-service?authSource=admin"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	ctx := context.Background()

	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", user, password, host, port, dbname)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	db := client.Database(dbname)

	partRepository, err := partRepository.NewRepository(db)
	if err != nil {
		log.Printf("failed to connect db: %v\n", err)
		return
	}

	service := partService.NewService(partRepository)

	api := partApiV1.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
