package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pptkna/rocket-factory/order/internal/config"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
	"go.uber.org/zap"
)

type app struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*app, error) {
	a := &app{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *app) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *app) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initOrderServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *app) initDI(_ context.Context) error {
	a.diContainer = newDIContainer()
	return nil
}

func (a *app) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *app) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *app) initOrderServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderApi.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderApi.ReadTimeout(),
		// Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	return nil
}

func (a *app) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP server listening on %s", config.AppConfig().OrderApi.Address()))

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, config.AppConfig().OrderApi.ShutDownTimeout())
		defer cancel()

		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "HTTP server crashed", zap.Error(err))

		return err
	}

	return nil
}
