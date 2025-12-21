module github.com/pptkna/rocket-factory/inventory

go 1.24.9

replace github.com/pptkna/rocket-factory/shared => ../shared

replace github.com/pptkna/rocket-factory/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/joho/godotenv v1.5.1
	github.com/pptkna/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.52.0
	go.mongodb.org/mongo-driver v1.17.6
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
