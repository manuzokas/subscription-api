package domain

import (
	"errors"
	"time"
)

// ErrSubscriptionNotFound é um erro customizado para quando uma assinatura não é encontrada.
var ErrSubscriptionNotFound = errors.New("subscription not found")

// ErrSubscriptionCannotBeCancelled é um erro para quando uma ação de cancelamento é inválida.
var ErrSubscriptionCannotBeCancelled = errors.New("subscription cannot be cancelled")

// ErrForbidden é usado quando um utilizador tenta aceder a um recurso que não lhe pertence.
var ErrForbidden = errors.New("user does not have permission to access this resource")

// Status define os possíveis estados de uma assinatura.
type Status string

const (
	StatusActive    Status = "ACTIVE"
	StatusTrial     Status = "TRIAL"
	StatusPastDue   Status = "PAST_DUE"
	StatusCancelled Status = "CANCELLED"
)

// Subscription é a entidade central do nosso domínio.
type Subscription struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	PlanID      string     `json:"planId"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CancelledAt *time.Time `json:"cancelledAt,omitempty"`
	TrialEndsAt *time.Time `json:"trialEndsAt,omitempty"`
}

// CanBeCancelled é um exemplo de regra de negócio dentro do domínio.
func (s *Subscription) CanBeCancelled() bool {
	return s.Status == StatusActive || s.Status == StatusTrial
}

// Cancel move a assinatura para o estado de cancelada.
func (s *Subscription) Cancel() {
	now := time.Now().UTC()
	s.Status = StatusCancelled
	s.CancelledAt = &now
	s.UpdatedAt = now
}
