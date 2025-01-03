package service

import (
	"context"
	"fmt"

	"github.com/emmrys-jay/ecommerce/internal/adapter/logger"
	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/emmrys-jay/ecommerce/internal/core/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

/**
 * OrderService implements port.OrderService interface
 */
type OrderService struct {
	repo        port.OrderRepository
	userRepo    port.UserRepository
	productRepo port.ProductRepository
	cache       port.CacheRepository
}

// NewOrderService creates a new order service instance
func NewOrderService(
	repo port.OrderRepository,
	userRepo port.UserRepository,
	productRepo port.ProductRepository,
	cache port.CacheRepository,
) *OrderService {
	return &OrderService{
		repo,
		userRepo,
		productRepo,
		cache,
	}
}

func (os *OrderService) PlaceOrder(ctx context.Context, userId uuid.UUID, req *domain.CreateOrderRequest) (*domain.Order, domain.CError) {
	_, cerr := os.userRepo.GetUserByID(ctx, userId)
	if cerr != nil {
		if cerr.Code() != 500 {
			return nil, domain.NewCError(cerr.Code(), "error fetching user: "+cerr.Error())
		}

		logger.FromCtx(ctx).Error("error fetching user to place order", zap.Error(cerr))
		return nil, domain.ErrInternal
	}

	var productIDMap = make(map[uuid.UUID]int)
	validProductIDs := make([]uuid.UUID, 0)
	for _, v := range req.Products {
		if id, err := uuid.Parse(v.ProductID); err == nil {
			validProductIDs = append(validProductIDs, id)
			productIDMap[id] = v.Quantity
		}
	}

	products, cerr := os.productRepo.GetProductsByIDs(ctx, validProductIDs)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error fetching products", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	if len(products) == 0 {
		return nil, domain.NewBadRequestCError("none of the products specified was found")
	}

	// Check for the integrity of order quantity with quantity in stock
	// Calculate total order amount
	// Populate order items
	var totalAmount float64
	var orderItems = make([]domain.OrderItem, 0)

	for _, v := range products {
		quantityOrdered := productIDMap[v.ID]

		if int(v.Quantity)-quantityOrdered < 0 {
			errMsg := fmt.Sprintf("The quantity specified for '%s' is more than the quantity in stock: %v (specified) for %v (in stock)",
				v.Name, productIDMap[v.ID], v.Quantity)
			return nil, domain.NewBadRequestCError(errMsg)
		}

		totalAmount += (v.Price * float64(quantityOrdered))

		orderItems = append(orderItems, domain.OrderItem{
			ProductID:   v.ID,
			ProductName: v.Name,
			Quantity:    int32(quantityOrdered),
			UnitPrice:   v.Price,
		})
	}

	order := domain.Order{
		UserID:      userId,
		TotalAmount: totalAmount,
		OrderItems:  orderItems,
		Status:      domain.OrderStatusPending,
	}

	retOrder, cerr := os.repo.CreateOrder(ctx, &order)
	if cerr != nil {
		if cerr.Code() == 500 {
			logger.FromCtx(ctx).Error("Error placing orders", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return retOrder, nil
}

func (os *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, domain.CError) {
	order, cerr := os.repo.GetOrder(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {
			logger.FromCtx(ctx).Error("Error getting single order", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return order, nil
}

func (os *OrderService) ListUserOrders(ctx context.Context, userId uuid.UUID) ([]domain.Order, domain.CError) {
	orders, cerr := os.repo.ListOrders(ctx, userId)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error fetching user orders", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return orders, nil
}

func (os *OrderService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, req *domain.UpdateOrderRequest) (*domain.Order, domain.CError) {
	retOrder, cerr := os.GetOrder(ctx, orderId)
	if cerr != nil {
		return nil, cerr
	}

	if _, ok := domain.StringToOrderStatus[req.Status]; !ok {
		return nil, domain.NewBadRequestCError("invalid status specified: " + req.Status)
	}

	order := domain.Order{
		ID:     orderId,
		Status: domain.StringToOrderStatus[req.Status],
	}

	_, cerr = os.repo.UpdateOrder(ctx, &order)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error updating order", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	retOrder.Status = order.Status // Add the updated status to the order struct to be returned
	return retOrder, nil
}

func (os *OrderService) CancelOrder(ctx context.Context, id uuid.UUID) (*domain.Order, domain.CError) {
	retOrder, cerr := os.GetOrder(ctx, id)
	if cerr != nil {
		return nil, cerr
	}

	if retOrder.Status != domain.OrderStatusPending {
		return nil, domain.NewBadRequestCError("You cannot cancel this order again since it has already been processed. Please contact admin")
	}

	order := domain.Order{
		ID:     id,
		Status: domain.OrderStatusCancelled,
	}

	_, cerr = os.repo.UpdateOrder(ctx, &order)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error canceling order", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	retOrder.Status = order.Status // Add the updated status to the order struct to be returned
	return retOrder, nil
}
