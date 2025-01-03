package port

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/google/uuid"
)

// ProductRepository is an interface for interacting with product-related data
type ProductRepository interface {
	// CreateProduct creates a new product and returns the created product
	CreateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError)
	// GetProductByID fetches a product specified by its id
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, domain.CError)
	// ListProducts fetches all products in the database
	ListProducts(ctx context.Context) ([]domain.Product, domain.CError)
	// UpdateProduct updates a product and returns the updated product
	UpdateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError)
	// DeleteProduct deletes a product specified by its id. It is a soft delete
	DeleteProduct(ctx context.Context, id uuid.UUID) domain.CError
	// GetProductsByIDs fetches all products that correspond to a list of product ids
	GetProductsByIDs(ctx context.Context, productIds []uuid.UUID) ([]domain.Product, domain.CError)
}

// ProductService is an interface for interacting with product-related business logic
type ProductService interface {
	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, prod *domain.CreateProductRequest) (*domain.Product, domain.CError)
	// GetProduct fetches a new product specified by its id
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, domain.CError)
	// ListProducts returns all products in the system
	ListProducts(ctx context.Context) ([]domain.Product, domain.CError)
	// UpdateProduct updates a products specified by its id
	UpdateProduct(ctx context.Context, id uuid.UUID, prod *domain.UpdateProductRequest) (*domain.Product, domain.CError)
	// DeleteProduct deletes a product in the system specified by its id
	DeleteProduct(ctx context.Context, id uuid.UUID) domain.CError
}
