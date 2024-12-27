package service

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/emmrys-jay/ecommerce/internal/core/port"
	"github.com/emmrys-jay/ecommerce/internal/core/util"
	"go.uber.org/zap"
)

/**
 * AuthService implements port.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	repo  port.UserRepository
	ts    port.TokenService
	cache port.CacheRepository
	l     *zap.Logger
}

// NewAuthService creates a new auth service instance
func NewAuthService(repo port.UserRepository, ts port.TokenService, cache port.CacheRepository, log *zap.Logger) *AuthService {
	return &AuthService{
		repo,
		ts,
		cache,
		log,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, domain.CError) {
	user, cerr := as.repo.GetUserByEmail(ctx, req.Email)
	if cerr != nil {
		if cerr == domain.ErrDataNotFound {
			return nil, domain.ErrInvalidCredentials
		}

		util.Error(as.l, ctx, "Error fetching user by email", cerr)
		return nil, domain.ErrInternal
	}

	err := util.ComparePassword(req.Password, user.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user.ID.String(), req.Email, user.Role.String())
	if err != nil {

		util.Error(as.l, ctx, "Error creating token", cerr)
		return nil, domain.ErrTokenCreation
	}
	util.Info(as.l, ctx, "created token for user using", "email/role", req.Email+"/"+user.Role.String())

	return &domain.LoginResponse{
		Token: accessToken,
		User:  *user,
	}, nil
}
