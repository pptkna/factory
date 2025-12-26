package app

import (
	"context"
	"errors"
	"net"

	"github.com/pptkna/rocket-factory/payment/internal/config"
	"github.com/pptkna/rocket-factory/platform/pkg/closer"
	"github.com/pptkna/rocket-factory/platform/pkg/grpc/health"
	"github.com/pptkna/rocket-factory/platform/pkg/logger"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type app struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
	return a.runGRPCServer(ctx)
}

func (a *app) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
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
	a.diContainer = newDiContainer()

	return nil
}

func (a *app) initLogger(_ context.Context) error {
	return logger.Init(config.AppConfig().Logger.Level(), config.AppConfig().Logger.AsJson())
}

func (a *app) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *app) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := a.listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *app) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Register health service for check availability
	health.RegisterService(a.grpcServer)

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1Api(ctx))

	return nil
}

func (a *app) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, "gRPC PaymentService server listening:", zap.String("address", config.AppConfig().PaymentGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
