package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manuzokas/subscription-api/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) Save(ctx context.Context, sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, plan_id, status, created_at, updated_at, cancelled_at, trial_ends_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			plan_id = EXCLUDED.plan_id,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at,
			cancelled_at = EXCLUDED.cancelled_at,
			trial_ends_at = EXCLUDED.trial_ends_at;
	`
	_, err := r.pool.Exec(ctx, query,
		sub.ID, sub.UserID, sub.PlanID, sub.Status,
		sub.CreatedAt, sub.UpdatedAt, sub.CancelledAt, sub.TrialEndsAt,
	)
	return err
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, created_at, updated_at, cancelled_at, trial_ends_at
		FROM subscriptions
		WHERE id = $1;
	`
	var sub domain.Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status,
		&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt, &sub.TrialEndsAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &sub, nil
}
