package auth

import (
	"context"

	"github.com/manuzokas/subscription-api/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
