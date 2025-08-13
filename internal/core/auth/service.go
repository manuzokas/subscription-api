package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/manuzokas/subscription-api/internal/domain"
)

var ErrUserAlreadyExists = errors.New("user with this email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	_, err := s.repo.FindUserByEmail(ctx, input.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.NewString(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*domain.User, error) {
	user, err := s.repo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
