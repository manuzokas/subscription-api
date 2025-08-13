package subscription

import (
	"context"

	"github.com/manuzokas/subscription-api/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, sub *domain.Subscription) error
	FindByID(ctx context.Context, id string) (*domain.Subscription, error)
}
