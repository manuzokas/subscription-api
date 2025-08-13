package subscription

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/manuzokas/subscription-api/internal/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type CreateSubscriptionInput struct {
	PlanID string `json:"planId" validate:"required"`
}

func (s *Service) CreateSubscription(ctx context.Context, userID string, input CreateSubscriptionInput) (*domain.Subscription, error) {
	now := time.Now().UTC()

	newSubscription := &domain.Subscription{
		ID:          uuid.NewString(),
		UserID:      userID,
		PlanID:      input.PlanID,
		Status:      domain.StatusTrial,
		CreatedAt:   now,
		UpdatedAt:   now,
		TrialEndsAt: &[]time.Time{now.Add(14 * 24 * time.Hour)}[0],
	}

	err := s.repo.Save(ctx, newSubscription)
	if err != nil {
		return nil, err
	}

	return newSubscription, nil
}

func (s *Service) GetSubscriptionByID(ctx context.Context, userID, subscriptionID string) (*domain.Subscription, error) {
	sub, err := s.repo.FindByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	if sub.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return sub, nil
}

func (s *Service) CancelSubscription(ctx context.Context, userID, subscriptionID string) error {
	sub, err := s.repo.FindByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if sub.UserID != userID {
		return domain.ErrForbidden
	}

	if !sub.CanBeCancelled() {
		return domain.ErrSubscriptionCannotBeCancelled
	}

	sub.Cancel()

	return s.repo.Save(ctx, sub)
}
