package port

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/google/uuid"
)

// OrderRepository is an interface for interacting with order-related data
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError)

	GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, domain.CError)

	ListOrders(ctx context.Context, userId uuid.UUID) ([]domain.Order, domain.CError)

	UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError)
}

// OrderService is an interface for interacting with order-related business logic
type OrderService interface {
	PlaceOrder(ctx context.Context, userId uuid.UUID, req *domain.CreateOrderRequest) (*domain.Order, domain.CError)

	GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, domain.CError)

	ListUserOrders(ctx context.Context, userId uuid.UUID) ([]domain.Order, domain.CError)

	UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, req *domain.UpdateOrderRequest) (*domain.Order, domain.CError)

	CancelOrder(ctx context.Context, id uuid.UUID) (*domain.Order, domain.CError)
}
