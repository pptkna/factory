package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	v1 "github.com/pptkna/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/pptkna/rocket-factory/order/internal/client/grpc"
	inventoryGRPCV1 "github.com/pptkna/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentGRPCV1 "github.com/pptkna/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/pptkna/rocket-factory/order/internal/config"
	kafkaConverter "github.com/pptkna/rocket-factory/order/internal/converter/kafka"
	decoder "github.com/pptkna/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/pptkna/rocket-factory/order/internal/repository"
	orderRepo "github.com/pptkna/rocket-factory/order/internal/repository/order"
	"github.com/pptkna/rocket-factory/order/internal/service"
	orderConsumer "github.com/pptkna/rocket-factory/order/internal/service/consumer/order_consumer"
	"github.com/pptkna/rocket-factory/order/internal/service/order"
	orderProducer "github.com/pptkna/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/pptkna/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/pptkna/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/pptkna/rocket-factory/platform/pkg/kafka/producer"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/pptkna/rocket-factory/platform/pkg/middleware/kafka"
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

	orderAssembledConsumerService service.OrderAssembledConsumerService
	orderPaidProducerService      service.OrderPaidProducerService

	orderAssembledConsumerGroup sarama.ConsumerGroup

	orderAssembledConsumer wrappedKafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.Producer

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
		d.orderService = order.NewService(d.OrderRepository(ctx), d.InventoryGRPCClient(ctx), d.PaymentGRPCClient(ctx), d.OrderPaidProducerService())
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

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderConsumer.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.OrderService(ctx),
		)
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidProducerService() service.OrderPaidProducerService {
	if d.orderPaidProducerService == nil {
		d.orderPaidProducerService = orderProducer.NewService(d.OrderPaidProducer())
	}

	return d.orderPaidProducerService
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		c, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.ConsumerGroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = c
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderPaidProducer() wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
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
