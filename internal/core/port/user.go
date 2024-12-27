package port

import (
	"context"

	"github.com/emmrys-jay/ecommerce/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository is an interface for interacting with User-related data
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, domain.CError)

	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, domain.CError)

	GetUserByEmail(ctx context.Context, email string) (*domain.User, domain.CError)

	ListUsers(ctx context.Context) ([]domain.User, domain.CError)

	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, domain.CError)

	DeleteUser(ctx context.Context, id uuid.UUID) domain.CError
}

// UserService is an interface for interacting with User-related business logic
type UserService interface {
	RegisterUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, domain.CError)

	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, domain.CError)

	ListUsers(ctx context.Context) ([]domain.User, domain.CError)

	UpdateUser(ctx context.Context, id uuid.UUID, user *domain.UpdateUserRequest) (*domain.User, domain.CError)

	DeleteUser(ctx context.Context, id uuid.UUID) domain.CError

	CreateAdminUser(ctx context.Context, email, password string) domain.CError
}
