package part

import (
	"sync"
	"time"

	def "github.com/pptkna/rocket-factory/inventory/internal/repository"
	repoModel "github.com/pptkna/rocket-factory/inventory/internal/repository/model"
	"github.com/samber/lo"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewPartRepository() *repository {
	return &repository{
		// parts: make(map[string]repoModel.Part),
		parts: map[string]repoModel.Part{
			"uuid-1": {
				Uuid:          "uuid-1",
				Name:          "Engine Turbocharger",
				Description:   "High-performance turbocharger for spacecraft engines.",
				Price:         1499.99,
				StockQuantity: 12,
				Category:      repoModel.CategoryEngine,
				Dimensions: &repoModel.Dimensions{
					Length: 42.5,
					Width:  28.1,
					Height: 19.7,
					Weight: 14.3,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    "NovaTech Industries",
					Country: "USA",
					Website: "https://novatech.space",
				},
				Tags: []string{"engine", "turbo", "performance"},
				Metadata: map[string]*repoModel.Value{
					"max_rpm":    {Int64: lo.ToPtr(int64(98000))},
					"efficiency": {Double: lo.ToPtr(float64(0.87))},
					"certified":  {Bool: lo.ToPtr(true)},
					"material":   {String: lo.ToPtr("Titanium Alloy")},
				},
				CreatedAt: lo.ToPtr(time.Date(2024, 10, 1, 12, 0, 0, 0, time.UTC)),
				UpdatedAt: lo.ToPtr(time.Date(2024, 12, 5, 16, 30, 0, 0, time.UTC)),
			},

			"uuid-2": {
				Uuid:          "uuid-2",
				Name:          "Fuel Injector",
				Description:   "Precision high-pressure injector for starship fuel systems.",
				Price:         320.75,
				StockQuantity: 40,
				Category:      repoModel.CategoryFuel,
				Dimensions: &repoModel.Dimensions{
					Length: 12.4,
					Width:  3.2,
					Height: 3.2,
					Weight: 0.8,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    "Galactic Fuel Systems",
					Country: "Germany",
					Website: "https://gfsystems.eu",
				},
				Tags: []string{"fuel", "injector", "high-pressure"},
				Metadata: map[string]*repoModel.Value{
					"pressure_rating": {Double: lo.ToPtr(float64(32))},
					"serial_number":   {String: lo.ToPtr("GFS-2024-11-2231")},
				},
				UpdatedAt: lo.ToPtr(time.Date(2024, 11, 2, 10, 10, 0, 0, time.UTC)),
			},
		},
	}
}
