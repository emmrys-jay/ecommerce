package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusActive     ProductStatus = "active"
	ProductStatusInactive   ProductStatus = "inactive"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
)

var StringToProductStatus = map[string]ProductStatus{
	"active":       ProductStatusActive,
	"inactive":     ProductStatusInactive,
	"out_of_stock": ProductStatusOutOfStock,
}

func (e *ProductStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProductStatus(s)
	case string:
		*e = ProductStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ProductStatus: %T", src)
	}
	return nil
}

func (ur ProductStatus) String() string {
	return string(ur)
}

type Product struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Quantity    int32         `json:"quantity"`
	Status      ProductStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   *time.Time    `json:"deleted_at"`
}

type CreateProductRequest struct {
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Price       float64       `json:"price" validate:"required,gte=0"`
	Quantity    int32         `json:"quantity" validate:"required,gte=1"`
	Status      ProductStatus `json:"-"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int32   `json:"quantity"`
	Status      string  `json:"status"`
}
