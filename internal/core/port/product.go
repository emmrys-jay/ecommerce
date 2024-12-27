package port

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/google/uuid"
)

// ProductRepository is an interface for interacting with product-related data
type ProductRepository interface {
	CreateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError)

	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, domain.CError)

	ListProducts(ctx context.Context) ([]domain.Product, domain.CError)

	UpdateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError)

	DeleteProduct(ctx context.Context, id uuid.UUID) domain.CError

	GetProductsByIDs(ctx context.Context, productIds []uuid.UUID) ([]domain.Product, domain.CError)
}

// ProductService is an interface for interacting with product-related business logic
type ProductService interface {
	CreateProduct(ctx context.Context, prod *domain.CreateProductRequest) (*domain.Product, domain.CError)

	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, domain.CError)

	ListProducts(ctx context.Context) ([]domain.Product, domain.CError)

	UpdateProduct(ctx context.Context, id uuid.UUID, prod *domain.UpdateProductRequest) (*domain.Product, domain.CError)

	DeleteProduct(ctx context.Context, id uuid.UUID) domain.CError
}
