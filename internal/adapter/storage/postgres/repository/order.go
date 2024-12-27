package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/emmrys-jay/ecommerce/internal/adapter/storage/postgres"
	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/**
 * OrderRepository implements port.OrderRepository interface
 * and provides an access to the postgres database
 */
type OrderRepository struct {
	db *postgres.DB
}

// NewOrderRepository creates a new order repository instance
func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError) {
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, domain.NewInternalCError("error starting transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	// Insert order first
	orderQuery := or.db.QueryBuilder.Insert("orders").
		Columns("user_id", "status", "total_amount").
		Values(order.UserID, order.Status, order.TotalAmount).
		Suffix("RETURNING *")

	sql, args, err := orderQuery.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError("error building query for orders: " + err.Error())
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, domain.NewInternalCError("error in scan for orders: " + err.Error())
	}

	// Build bulk insert for order items
	oiQuery := or.db.QueryBuilder.Insert("order_items").
		Columns("order_id", "product_id", "product_name", "quantity", "unit_price")

	for _, item := range order.OrderItems {
		oiQuery = oiQuery.Values(
			order.ID,
			item.ProductID,
			item.ProductName,
			item.Quantity,
			item.UnitPrice,
		)
	}
	oiQuery = oiQuery.Suffix("RETURNING id, created_at")

	sql, args, err = oiQuery.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError("error building query for order items: " + err.Error())
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		// 23503 is the error code for a foreign key violation error
		if errCode := or.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, domain.NewInternalCError("error executing query for order_items: " + err.Error())
	}
	// transaction won't commit if you defer closing row here
	// defer rows.Close()

	for i := range order.OrderItems {
		if !rows.Next() {
			return nil, domain.NewInternalCError("unexpected end of returned rows for order_items")
		}

		if err := rows.Scan(&order.OrderItems[i].ID, &order.OrderItems[i].CreatedAt); err != nil {
			return nil, domain.NewInternalCError("error getting returned rows for order_items: " + err.Error())
		}
	}

	if err = rows.Err(); err != nil {
		// 23503 is the error code for a foreign key violation error
		if errCode := or.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, domain.NewInternalCError("error getting returned rows for order_items: " + err.Error())
	}
	rows.Close() // Close row manually so the transaction can be committed

	err = tx.Commit(ctx)
	if err != nil {
		return nil, domain.NewInternalCError("error commiting transactions: " + err.Error())
	}

	return order, nil
}

func (or *OrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, domain.CError) {
	query := or.db.QueryBuilder.Select(
		"o.id",
		"o.user_id",
		"o.status",
		"o.total_amount",
		"o.created_at",
		"o.updated_at",
		"oi.id as item_id",
		"oi.order_id",
		"oi.product_id",
		"oi.product_name",
		"oi.quantity",
		"oi.unit_price",
		"oi.created_at as item_created_at",
	).
		From("orders o").
		Join("order_items oi ON o.id = oi.order_id").
		Where(sq.Eq{"o.id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError("error building sql: " + err.Error())
	}

	rows, err := or.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer rows.Close()

	var order = &domain.Order{}
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,

			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, domain.NewInternalCError(err.Error())
		}

		order.OrderItems = append(order.OrderItems, item)
	}

	if err := rows.Err(); err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(err.Error())
	}

	return order, nil
}

func (or *OrderRepository) ListOrders(ctx context.Context, userId uuid.UUID) ([]domain.Order, domain.CError) {
	query := or.db.QueryBuilder.Select(
		"o.id",
		"o.user_id",
		"o.status",
		"o.total_amount",
		"o.created_at",
		"o.updated_at",
		"oi.id as item_id",
		"oi.order_id",
		"oi.product_id",
		"oi.product_name",
		"oi.quantity",
		"oi.unit_price",
		"oi.created_at as item_created_at",
	).
		From("orders o").
		Join("order_items oi ON o.id = oi.order_id").
		Where(sq.Eq{"o.user_id": userId}).
		OrderBy("o.created_at DESC") // Most recent orders first

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError("error building sql: " + err.Error())
	}

	rows, err := or.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer rows.Close()

	// Map to store orders by ID to avoid duplicates
	orderMap := make(map[uuid.UUID]*domain.Order)
	orders := make([]domain.Order, 0)

	for rows.Next() {
		var (
			order domain.Order
			item  domain.OrderItem
		)

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,

			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.UnitPrice,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, domain.NewInternalCError("error scanning row: " + err.Error())
		}

		// Check if we've already seen this order
		existingOrder, exists := orderMap[order.ID]
		if !exists {
			// New order, initialize its items slice
			order.OrderItems = make([]domain.OrderItem, 0)
			orderMap[order.ID] = &order
			orders = append(orders, order)
			existingOrder = &orders[len(orders)-1]
		}

		existingOrder.OrderItems = append(existingOrder.OrderItems, item)
	}

	return orders, nil
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError) {
	query := or.db.QueryBuilder.Update("orders").
		Set("status", order.Status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": order.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	err = or.db.QueryRow(ctx, sql, args...).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(err.Error())
	}

	return order, nil
}
