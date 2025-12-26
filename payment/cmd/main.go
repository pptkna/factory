package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/pptkna/rocket-factory/payment/internal/app"
	"github.com/pptkna/rocket-factory/payment/internal/config"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to create new app", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to start app", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "failed to gracefull shutdown", zap.Error(err))
	}
}
