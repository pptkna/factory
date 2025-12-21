package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pptkna/rocket-factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	OrderApi      ApiConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderApiCfg, err := env.NewOrderApiConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		OrderApi:      orderApiCfg,
		Postgres:      postgresCfg,
		InventoryGRPC: inventoryGRPCCfg,
		PaymentGRPC:   paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
