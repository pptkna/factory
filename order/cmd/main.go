package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	orderApiV1 "github.com/pptkna/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/pptkna/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/pptkna/rocket-factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/pptkna/rocket-factory/order/internal/repository/order"
	orderService "github.com/pptkna/rocket-factory/order/internal/service/order"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	// Services
	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"

	// DB
	host     = "localhost"
	port     = "5432"
	user     = "order-service-user"
	password = "order-service-password"
	dbname   = "order-service"
	sslmode  = "disable"

	migrations_dir = "migrations"
)

func main() {
	con, err := orderRepository.NewRepository(host, port, user, password, dbname, sslmode, migrations_dir)
	if err != nil {
		log.Printf("failed to connect db: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(inventoryServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to inventory connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close inventory connect: %v", cerr)
		}
	}()

	inventoryServiceClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	inventoryClient := inventoryClient.NewClient(inventoryServiceClient)

	paymentConn, err := grpc.NewClient(paymentServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to payment connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close payment connect: %v", cerr)
		}
	}()

	paymentServerClient := paymentV1.NewPaymentServiceClient(paymentConn)

	paymentClient := paymentClient.NewClient(paymentServerClient)

	service := orderService.NewService(con, inventoryClient, paymentClient)

	api := orderApiV1.NewApi(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		// Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
