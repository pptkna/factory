module github.com/pptkna/rocket-factory/assembly

go 1.24.9

replace github.com/pptkna/rocket-factory/shared => ../shared

replace github.com/pptkna/rocket-factory/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/joho/godotenv v1.5.1
	github.com/pptkna/rocket-factory/shared v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.11
)
