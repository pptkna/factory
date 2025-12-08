package part

import (
	"context"
	"fmt"
	"time"

	def "github.com/pptkna/rocket-factory/inventory/internal/repository"
	"github.com/pptkna/rocket-factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) (*repository, error) {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "uuid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "category", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "manufacturer.country", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys: bson.D{
				{Key: "tags", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return nil, fmt.Errorf("failed to create database indexes: %w", err)
	}

	SeedTestParts(ctx, collection)

	return &repository{
		collection: collection,
	}, nil
}

// Add test data
func SeedTestParts(ctx context.Context, collection *mongo.Collection) error {
	now := time.Now()

	parts := []model.Part{
		{
			Uuid:          "c2a4e5cd-8aa1-4dd8-ab19-5cdf25af047d",
			Name:          "Turbo Engine X1",
			Description:   "High-performance turbo engine",
			Price:         12999.99,
			StockQuantity: 10,
			Category:      model.CategoryEngine,
			Dimensions: &model.Dimensions{
				Length: 120,
				Width:  80,
				Height: 70,
				Weight: 350,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "AeroTech",
				Country: "USA",
				Website: "https://aerotech.example.com",
			},
			Tags:      []string{"engine", "turbo", "performance"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "42d9bd89-6023-4d95-af31-c1d86bb39b43",
			Name:          "Fuel Filter S9",
			Description:   "Advanced fuel filtration system",
			Price:         349.50,
			StockQuantity: 120,
			Category:      model.CategoryFuel,
			Dimensions: &model.Dimensions{
				Length: 20,
				Width:  15,
				Height: 15,
				Weight: 3,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "FuelMaster",
				Country: "Germany",
				Website: "https://fuelmaster.example.com",
			},
			Tags:      []string{"fuel", "filter"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "19b5cae8-54d3-4fc8-97fa-bab9d0718d80",
			Name:          "Wing Panel A7",
			Description:   "Reinforced composite wing panel",
			Price:         5599.00,
			StockQuantity: 25,
			Category:      model.CategoryWing,
			Dimensions: &model.Dimensions{
				Length: 450,
				Width:  200,
				Height: 15,
				Weight: 220,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "SkyParts",
				Country: "Japan",
				Website: "https://skyparts.example.com",
			},
			Tags:      []string{"wing", "panel", "composite"},
			Metadata:  map[string]*model.Value{},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	for _, part := range parts {
		_, err := collection.InsertOne(ctx, part)
		if err != nil {
			return fmt.Errorf("failed to insert part %s: %w", part.Uuid, err)
		}
	}

	return nil
}
