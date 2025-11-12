package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	inventory_v1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

// InventoryService реализует gRPC сервис для работы с деталями
type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventory_v1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uuid := req.GetUuid()

	part, ok := s.parts[uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	s.mu.RLock()

	parts := make([]*inventory_v1.Part, 0, len(s.parts))

	for _, p := range s.parts {
		parts = append(parts, p)
	}

	s.mu.RUnlock()

	filters := req.GetFilter()
	if len(filters.Uuids) == 0 && len(filters.Categories) == 0 && len(filters.ManufacturerCountries) == 0 && len(filters.Names) == 0 && len(filters.Tags) == 0 {
		return &inventory_v1.ListPartsResponse{Parts: parts}, nil
	}

	filteredParts := make([]*inventory_v1.Part, 0, len(parts))
	for _, p := range parts {
		if matchesFilter(p, filters) {
			filteredParts = append(filteredParts, p)
		}
	}

	return &inventory_v1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

func matchesFilter(part *inventory_v1.Part, filters *inventory_v1.PartsFilter) bool {
	if filters == nil {
		return true
	}

	if uuids := filters.GetUuids(); len(uuids) > 0 {
		if !slices.Contains(uuids, part.Uuid) {
			return false
		}
	}

	if categories := filters.GetCategories(); len(categories) > 0 {
		if !slices.Contains(categories, part.Category) {
			return false
		}
	}

	if manufacturerCountries := filters.GetManufacturerCountries(); len(manufacturerCountries) > 0 {
		if !slices.Contains(manufacturerCountries, part.Manufacturer.Country) {
			return false
		}
	}

	if names := filters.GetNames(); len(names) > 0 {
		if !slices.Contains(names, part.Name) {
			return false
		}
	}

	if tags := filters.GetTags(); len(tags) > 0 {
		found := false

		for _, t := range tags {
			if slices.Contains(part.Tags, t) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{
		parts: make(map[string]*inventory_v1.Part),
	}

	inventory_v1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
