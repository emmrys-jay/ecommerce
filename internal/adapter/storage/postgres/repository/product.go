package repository

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/emmrys-jay/ecommerce/internal/adapter/storage/postgres"
	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/**
 * ProductRepository implements port.ProductRepository interface
 * and provides an access to the postgres database
 */
type ProductRepository struct {
	db *postgres.DB
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *postgres.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

func (ur *ProductRepository) CreateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError) {
	query := ur.db.QueryBuilder.Insert("products").
		Columns("name", "description", "price", "quantity", "status").
		Values(prod.Name, prod.Description, prod.Price, prod.Quantity, prod.Status).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	log.Println(sql)

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&prod.ID,
		&prod.Name,
		&prod.Description,
		&prod.Price,
		&prod.Quantity,
		&prod.Status,
		&prod.CreatedAt,
		&prod.UpdatedAt,
		&prod.DeletedAt,
	)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return prod, nil
}

// GetProductByID gets a product by its ID from the database
func (ur *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, domain.CError) {
	var prod domain.Product

	query := ur.db.QueryBuilder.Select("*").
		From("products").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&prod.ID,
		&prod.Name,
		&prod.Description,
		&prod.Price,
		&prod.Quantity,
		&prod.Status,
		&prod.CreatedAt,
		&prod.UpdatedAt,
		&prod.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(err.Error())
	}

	return &prod, nil
}

// ListProducts lists all products in the database
func (ur *ProductRepository) ListProducts(ctx context.Context) ([]domain.Product, domain.CError) {
	var prod domain.Product
	var prods []domain.Product

	query := ur.db.QueryBuilder.Select("*").
		From("products").
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Description,
			&prod.Price,
			&prod.Quantity,
			&prod.Status,
			&prod.CreatedAt,
			&prod.UpdatedAt,
			&prod.DeletedAt,
		)
		if err != nil {
			return nil, domain.NewInternalCError(err.Error())
		}

		prods = append(prods, prod)
	}

	return prods, nil
}

// UpdateProduct updates a product by ID in the database
func (ur *ProductRepository) UpdateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError) {
	query := ur.db.QueryBuilder.Update("products").
		Set("name", prod.Name).
		Set("description", prod.Description).
		Set("price", prod.Price).
		Set("quantity", prod.Quantity).
		Set("status", prod.Status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": prod.ID, "deleted_at": nil}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&prod.ID,
		&prod.Name,
		&prod.Description,
		&prod.Price,
		&prod.Quantity,
		&prod.Status,
		&prod.CreatedAt,
		&prod.UpdatedAt,
		&prod.DeletedAt,
	)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return prod, nil
}

// DeleteProduct deletes a product by ID from the database
func (ur *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) domain.CError {
	query := ur.db.QueryBuilder.Update("products").
		Set("deleted_at", time.Now()).
		Set("status", domain.ProductStatusInactive).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.NewInternalCError(err.Error())
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrDataNotFound
		}
		return domain.NewInternalCError(err.Error())
	}

	return nil
}

// GetProductsByIDs gets a number of products by their ids
func (ur *ProductRepository) GetProductsByIDs(ctx context.Context, productIds []uuid.UUID) ([]domain.Product, domain.CError) {
	var prod domain.Product
	var prods []domain.Product

	query := ur.db.QueryBuilder.Select("*").
		From("products").
		Where(sq.Eq{"deleted_at": nil, "id": productIds, "status": domain.ProductStatusActive}).
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Description,
			&prod.Price,
			&prod.Quantity,
			&prod.Status,
			&prod.CreatedAt,
			&prod.UpdatedAt,
			&prod.DeletedAt,
		)
		if err != nil {
			return nil, domain.NewInternalCError(err.Error())
		}

		prods = append(prods, prod)
	}

	return prods, nil
}
