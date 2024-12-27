package service

import (
	"context"
	"net/http"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/emmrys-jay/ecommerce/internal/core/port"
	"github.com/emmrys-jay/ecommerce/internal/core/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

/**
 * ProductService implements port.ProductService interface
 */
type ProductService struct {
	repo  port.ProductRepository
	cache port.CacheRepository
	l     *zap.Logger
}

// NewProductService creates a new product service instance
func NewProductService(repo port.ProductRepository, cache port.CacheRepository, log *zap.Logger) *ProductService {
	return &ProductService{
		repo,
		cache,
		log,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, prod *domain.CreateProductRequest) (*domain.Product, domain.CError) {
	prodToCreate := domain.Product{
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		Quantity:    prod.Quantity,
		Status:      domain.ProductStatusActive,
	}

	prodResponse, cerr := ps.repo.CreateProduct(ctx, &prodToCreate)
	if cerr != nil {
		if cerr.Code() == 500 {

			util.Error(ps.l, ctx, "Error creating product", cerr)
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return prodResponse, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, domain.CError) {
	product, cerr := ps.repo.GetProductByID(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {

			util.Error(ps.l, ctx, "Error getting product", cerr)
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return product, nil
}

func (ps *ProductService) ListProducts(ctx context.Context) ([]domain.Product, domain.CError) {
	users, cerr := ps.repo.ListProducts(ctx)
	if cerr != nil {

		util.Error(ps.l, ctx, "Error listing products", cerr)
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req *domain.UpdateProductRequest) (*domain.Product, domain.CError) {
	retProd, cerr := ps.GetProduct(ctx, id)
	if cerr != nil {
		return nil, cerr
	}

	if req.Name == retProd.Name && req.Description == retProd.Description && req.Status == retProd.Status.String() &&
		req.Price == retProd.Price && req.Quantity == retProd.Quantity {
		return nil, domain.NewCError(http.StatusBadRequest, "There are no changes to update")
	}

	retProd.Name = req.Name
	retProd.Description = req.Description
	retProd.Price = req.Price
	retProd.Quantity = req.Quantity

	if status, ok := domain.StringToProductStatus[req.Status]; ok {
		retProd.Status = status
	}

	userResponse, cerr := ps.repo.UpdateProduct(ctx, retProd)
	if cerr != nil {
		if cerr.Code() == 500 {

			util.Error(ps.l, ctx, "Error updating product", cerr)
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return userResponse, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) domain.CError {
	cerr := ps.repo.DeleteProduct(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {

			util.Error(ps.l, ctx, "Error deleting product", cerr)
			return domain.ErrInternal
		}
		return cerr
	}

	return nil
}
