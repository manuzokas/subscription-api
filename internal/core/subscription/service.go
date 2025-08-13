package subscription

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/manuzokas/subscription-api/internal/core/auth"
	"github.com/manuzokas/subscription-api/internal/domain"
)

type MessagePublisher interface {
	Publish(ctx context.Context, queueName string, body []byte) error
}

type Service struct {
	repo      Repository
	userRepo  auth.UserRepository
	publisher MessagePublisher
}

func NewService(repo Repository, userRepo auth.UserRepository, publisher MessagePublisher) *Service {
	return &Service{
		repo:      repo,
		userRepo:  userRepo,
		publisher: publisher,
	}
}

type CreateSubscriptionInput struct {
	PlanID string `json:"planId" validate:"required"`
}

type SubscriptionCreatedEvent struct {
	SubscriptionID string `json:"subscriptionId"`
	UserID         string `json:"userId"`
	Email          string `json:"email"`
}

func (s *Service) CreateSubscription(ctx context.Context, userID string, input CreateSubscriptionInput) (*domain.Subscription, error) {
	now := time.Now().UTC()

	newSubscription := &domain.Subscription{
		ID:        uuid.NewString(),
		UserID:    userID,
		PlanID:    input.PlanID,
		Status:    "PENDING",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Save(ctx, newSubscription); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	event := SubscriptionCreatedEvent{
		SubscriptionID: newSubscription.ID,
		UserID:         newSubscription.UserID,
		Email:          user.Email,
	}
	eventBody, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	const queueName = "subscription_created_events"
	if err := s.publisher.Publish(ctx, queueName, eventBody); err != nil {
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
