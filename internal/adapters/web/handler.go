package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/manuzokas/subscription-api/internal/core/subscription"
	"github.com/manuzokas/subscription-api/internal/domain"
)

type SubscriptionHandler struct {
	service *subscription.Service
}

func NewSubscriptionHandler(s *subscription.Service) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: s,
	}
}

func (h *SubscriptionHandler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok {
		http.Error(w, "invalid user ID in context", http.StatusInternalServerError)
		return
	}

	var input subscription.CreateSubscriptionInput
	if err := decodeAndValidate(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := h.service.CreateSubscription(r.Context(), userID, input)
	if err != nil {
		http.Error(w, "could not create subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) GetSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok {
		http.Error(w, "invalid user ID in context", http.StatusInternalServerError)
		return
	}
	subscriptionID := chi.URLParam(r, "id")

	sub, err := h.service.GetSubscriptionByID(r.Context(), userID, subscriptionID)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "could not retrieve subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) CancelSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok {
		http.Error(w, "invalid user ID in context", http.StatusInternalServerError)
		return
	}
	subscriptionID := chi.URLParam(r, "id")

	err := h.service.CancelSubscription(r.Context(), userID, subscriptionID)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if errors.Is(err, domain.ErrSubscriptionCannotBeCancelled) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "could not cancel subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
