package service

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/emmrys-jay/ecommerce/internal/core/port"
	"go.uber.org/zap"
)

/**
 * PingService implements port.PingService interface
 */
type PingService struct {
	repo  port.PingRepository
	cache port.CacheRepository
	l     *zap.Logger
}

// NewAuthService creates a new auth service instance
func NewPingService(repo port.PingRepository, cache port.CacheRepository, log *zap.Logger) *PingService {
	return &PingService{
		repo,
		cache,
		log,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (ps *PingService) Ping(ctx context.Context, ping *domain.Ping) (domain.Ping, domain.CError) {
	_ = ps.repo.CreatePing(ctx, ping)
	return *ping, nil
}
