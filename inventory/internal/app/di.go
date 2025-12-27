package app

import (
	"context"
	"fmt"

	v1 "github.com/pptkna/rocket-factory/inventory/internal/api/part/v1"
	"github.com/pptkna/rocket-factory/inventory/internal/config"
	"github.com/pptkna/rocket-factory/inventory/internal/repository"
	partRepo "github.com/pptkna/rocket-factory/inventory/internal/repository/part"
	"github.com/pptkna/rocket-factory/inventory/internal/service"
	partService "github.com/pptkna/rocket-factory/inventory/internal/service/part"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	partService service.PartService

	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDB       *mongo.Database
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = v1.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		partRepository, err := partRepo.NewRepository(d.MongoDB(ctx))
		if err != nil {
			panic(fmt.Sprintf("failed to create part repository: %s\n", err.Error()))
		}

		d.partRepository = partRepository
	}

	return d.partRepository
}

func (d *diContainer) MongoDB(ctx context.Context) *mongo.Database {
	if d.mongoDB == nil {
		db := d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
		d.mongoDB = db
	}

	return d.mongoDB
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to create mongoDB client: %s\n", err.Error()))
		}
		closer.AddNamed("mongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		err = client.Ping(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to ping database: %s\n", err.Error()))
		}

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}
