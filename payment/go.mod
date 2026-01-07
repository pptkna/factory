module github.com/pptkna/rocket-factory/payment

go 1.24.9

replace github.com/pptkna/rocket-factory/shared => ../shared

replace github.com/pptkna/rocket-factory/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/pptkna/rocket-factory/platform v0.0.0-00010101000000-000000000000
	github.com/pptkna/rocket-factory/shared v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.77.0
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
