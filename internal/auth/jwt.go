package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hydrocode-de/gorun/internal/db"
)

func CreateJWT(refreshToken string, secretKey string, validFor time.Duration, DB *db.Queries, ctx context.Context) (string, error) {
	user, err := DB.GetRefreshTokenUser(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(validFor).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, secretKey string) (string, error) {
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id claim not found")
	}

	return userId, nil
}
