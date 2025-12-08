package model

import "time"

type Category string

const (
	CategoryUnspecified Category = "UNSPECIFIED"
	CategoryUnknown     Category = "UNKNOWN"
	CategoryEngine      Category = "ENGINE"
	CategoryFuel        Category = "FUEL"
	CategoryPorthole    Category = "PORTHOLE"
	CategoryWing        Category = "WING"
)

type PartFiters struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Dimensions struct {
	Length float32 `bson:"length"`
	Width  float32 `bson:"width"`
	Height float32 `bson:"height"`
	Weight float32 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}
type Value struct {
	String *string  `bson:"string,omitempty"`
	Int64  *int64   `bson:"int64,omitempty"`
	Double *float64 `bson:"double,omitempty"`
	Bool   *bool    `bson:"bool,omitempty"`
}

type Part struct {
	Uuid          string            `bson:"uuid"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	Price         float32           `bson:"price"`
	StockQuantity int               `bson:"stock_quantity"`
	Category      Category          `bson:"category"`
	Dimensions    *Dimensions       `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer     `bson:"manufacturer,omitempty"`
	Tags          []string          `bson:"tags,omitempty"`
	Metadata      map[string]*Value `bson:"metadata,omitempty"`
	CreatedAt     *time.Time        `bson:"created_at,omitempty"`
	UpdatedAt     *time.Time        `bson:"updated_at,omitempty"`
}
