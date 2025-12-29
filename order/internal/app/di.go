package app

import (
	"context"
	"fmt"

	v1 "github.com/pptkna/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/pptkna/rocket-factory/order/internal/client/grpc"
	inventoryGRPCV1 "github.com/pptkna/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentGRPCV1 "github.com/pptkna/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/pptkna/rocket-factory/order/internal/config"
	"github.com/pptkna/rocket-factory/order/internal/repository"
	orderRepo "github.com/pptkna/rocket-factory/order/internal/repository/order"
	"github.com/pptkna/rocket-factory/order/internal/service"
	"github.com/pptkna/rocket-factory/order/internal/service/order"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService service.OrderService

	inventoryGRPCClient grpcClient.InventoryClient
	paymentGRPCClient   grpcClient.PaymentClient

	orderRepository repository.OrderRepository
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = v1.NewApi(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = order.NewService(d.OrderRepository(ctx), d.InventoryGRPCClient(ctx), d.PaymentGRPCClient(ctx))
	}

	return d.orderService
}

func (d *diContainer) InventoryGRPCClient(ctx context.Context) grpcClient.InventoryClient {
	if d.inventoryGRPCClient == nil {
		inventoryConn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("failed to inventory client connect: %s\n", err.Error()))
		}
		closer.AddNamed("inventory client connection", func(ctx context.Context) error {
			return inventoryConn.Close()
		})

		inventoryServiceClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

		d.inventoryGRPCClient = inventoryGRPCV1.NewClient(inventoryServiceClient)
	}

	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentGRPCClient(ctx context.Context) grpcClient.PaymentClient {
	if d.paymentGRPCClient == nil {
		paymentConn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("failed to payment client connect: %s\n", err.Error()))
		}
		closer.AddNamed("inventory client connection", func(ctx context.Context) error {
			return paymentConn.Close()
		})

		paymentServerClient := paymentV1.NewPaymentServiceClient(paymentConn)

		d.paymentGRPCClient = paymentGRPCV1.NewClient(paymentServerClient)
	}

	return d.paymentGRPCClient
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		con, err := orderRepo.NewRepository(config.AppConfig().Postgres.Address(), config.AppConfig().Postgres.MigrationDirectory())
		if err != nil {
			panic(fmt.Sprintf("failed to connect db: %s\n", err.Error()))
		}

		closer.AddNamed("order repository", func(ctx context.Context) error {
			return con.Close()
		})

		d.orderRepository = con
	}

	return d.orderRepository
}
