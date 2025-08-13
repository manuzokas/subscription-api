package web

import (
	"encoding/json"
	"net/http"

	"github.com/manuzokas/subscription-api/internal/core/auth"
)

type AuthHandler struct {
	service   *auth.AuthService
	jwtSecret string
}

func NewAuthHandler(s *auth.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		service:   s,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input auth.RegisterInput
	if err := decodeAndValidate(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input auth.LoginInput
	if err := decodeAndValidate(r, &input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.ID, h.jwtSecret)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
