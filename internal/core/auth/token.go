package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}
