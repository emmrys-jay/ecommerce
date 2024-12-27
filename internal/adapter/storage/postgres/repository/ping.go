package repository

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/adapter/storage/postgres"
	"github.com/emmrys-jay/ecommerce/internal/core/domain"
)

/**
 * CategoryRepository implements port.CategoryRepository interface
 * and provides an access to the postgres database
 */
type PingRepository struct {
	db *postgres.DB
}

// NewCategoryRepository creates a new category repository instance
func NewPingRepository(db *postgres.DB) *PingRepository {
	return &PingRepository{
		db,
	}
}

func (pr *PingRepository) CreatePing(ctx context.Context, category *domain.Ping) error {
	return nil
}
