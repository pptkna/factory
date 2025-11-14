package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	orderV1 "github.com/pptkna/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
	payment_v1 "github.com/pptkna/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
)

var (
	notFoundError = errors.New("not found")
	conflictError = errors.New("conflict error")
)

type orderStorage struct {
	mu           sync.RWMutex
	orders       map[string]*orderV1.OrderDto
	lockedOrders map[string]struct{}
}

func NewOrderStorage() *orderStorage {
	return &orderStorage{
		orders:       make(map[string]*orderV1.OrderDto),
		lockedOrders: make(map[string]struct{}),
	}
}

// Создание заказа
func (s *orderStorage) CreateOrder(uuid string, order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.orders[uuid]; exist {
		return conflictError
	}

	s.orders[uuid] = order

	return nil
}

// Получение информации о заказе
func (s *orderStorage) GetOrder(uuid string) (*orderV1.OrderDto, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, notFoundError
	}

	copy := *order
	return &copy, nil
}

// Блокировка заказа перед обновлением
func (s *orderStorage) LockOrder(uuid string) (*orderV1.OrderDto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, notFoundError
	}

	if _, locked := s.lockedOrders[uuid]; locked {
		return nil, conflictError
	}

	s.lockedOrders[uuid] = struct{}{}

	copy := *order

	return &copy, nil
}

// Разблокировка заказа, если обновление не требуется
func (s *orderStorage) UnlockOrder(uuid string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.lockedOrders, uuid)
}

// Обнвление заблокированного заказа и последующая разблокировка
func (s *orderStorage) UpdateLockedOrder(uuid string, update func(o *orderV1.OrderDto)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, locked := s.lockedOrders[uuid]; !locked {
		return conflictError
	}

	order, ok := s.orders[uuid]
	if !ok {
		delete(s.lockedOrders, uuid)
		return notFoundError
	}

	defer delete(s.lockedOrders, uuid)

	update(order)

	return nil
}

type orderHandler struct {
	storage                *orderStorage
	inventoryServiceClient inventoryV1.InventoryServiceClient
	paymentServiceClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *orderStorage, inventoryServiceClient inventoryV1.InventoryServiceClient, paymentServiceClient paymentV1.PaymentServiceClient) *orderHandler {
	return &orderHandler{
		storage,
		inventoryServiceClient,
		paymentServiceClient,
	}
}

func (h *orderHandler) GetOrderByOrderUuid(_ context.Context, params orderV1.GetOrderByOrderUuidParams) (orderV1.GetOrderByOrderUuidRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID)
	if err != nil {
		if errors.Is(err, notFoundError) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with uuid: %s not found", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return order, nil
}

func (h *orderHandler) PostOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.PostOrderRes, error) {
	parts, err := h.inventoryServiceClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: mapUuidToStr(req.GetPartUuids()),
		},
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	if len(parts.GetParts()) != len(req.GetPartUuids()) {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Some parts not found",
		}, nil
	}

	var totalPrice float32

	for _, p := range parts.GetParts() {
		if p == nil {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: "Internal server error",
			}, nil
		}

		totalPrice += float32(p.GetPrice())
	}

	uuid := uuid.New()

	err = h.storage.CreateOrder(uuid.String(), &orderV1.OrderDto{
		OrderUUID:  uuid,
		UserUUID:   req.GetUserUUID(),
		PartUuids:  req.GetPartUuids(),
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return &orderV1.CreateOrderResponse{}, nil
}

func (h *orderHandler) PostOrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PostOrderPayParams) (orderV1.PostOrderPayRes, error) {
	order, err := h.storage.LockOrder(params.OrderUUID)
	defer h.storage.UnlockOrder(params.OrderUUID)
	if err != nil {
		if errors.Is(err, conflictError) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Order already in payment process",
			}, nil
		}
		if errors.Is(err, notFoundError) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with uuid: %s not found", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order already paid",
		}, nil
	}

	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order already cancelled",
		}, nil
	}

	payment, err := h.paymentServiceClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID,
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: mapPaymentMethodFromOrderToPayment(req.PaymentMethod),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	h.storage.UpdateLockedOrder(params.OrderUUID, func(order *orderV1.OrderDto) {
		order.Status = orderV1.OrderStatusPAID
	})

	transactionUUID, err := uuid.Parse(payment.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (h *orderHandler) PostOrderCancel(ctx context.Context, params orderV1.PostOrderCancelParams) (orderV1.PostOrderCancelRes, error) {
	order, err := h.storage.LockOrder(params.OrderUUID)
	defer h.storage.UnlockOrder(params.OrderUUID)
	if err != nil {
		if errors.Is(err, conflictError) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Order already in payment process",
			}, nil
		}
		if errors.Is(err, notFoundError) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order with uuid: %s not found", params.OrderUUID),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order already paid",
		}, nil
	}

	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order already cancelled",
		}, nil
	}

	h.storage.UpdateLockedOrder(params.OrderUUID, func(order *orderV1.OrderDto) {
		order.Status = orderV1.OrderStatusCANCELLED
	})

	return &orderV1.PostOrderCancelNoContent{}, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *orderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func mapUuidToStr(uuids []uuid.UUID) []string {
	strUuids := make([]string, len(uuids))

	for i, uuid := range uuids {
		strUuids[i] = uuid.String()
	}

	return strUuids
}

func mapPaymentMethodFromOrderToPayment(pm orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch pm {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case orderV1.PaymentMethodUNKNOWN:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func main() {
	storage := NewOrderStorage()

	inventoryConn, err := grpc.NewClient(inventoryServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to inventory connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close inventory connect: %v", cerr)
		}
	}()

	inventoryServerClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(paymentServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to payment connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close payment connect: %v", cerr)
		}
	}()

	paymentServerClient := paymentV1.NewPaymentServiceClient(paymentConn)

	orderHandler := NewOrderHandler(storage, inventoryServerClient, paymentServerClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		// Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
