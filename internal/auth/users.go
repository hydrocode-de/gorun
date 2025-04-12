package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hydrocode-de/gorun/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginResponse struct {
	User         User      `json:"user"`
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token"`
}

func CreateUser(ctx context.Context, DB *db.Queries, email, password string, isAdmin bool, jwtSecret string) (UserLoginResponse, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserLoginResponse{}, err
	}

	user, err := DB.CreateUser(ctx, db.CreateUserParams{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPw),
		IsAdmin:      isAdmin,
	})
	if err != nil {
		return UserLoginResponse{}, err
	}

	refreshToken := uuid.New().String()
	_, err = DB.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	})
	if err != nil {
		return UserLoginResponse{}, err
	}

	return NewJWTFromRefreshToken(ctx, DB, refreshToken, jwtSecret)
}

func LoginUser(ctx context.Context, DB *db.Queries, email, password, jwtSecret string) (UserLoginResponse, error) {
	user, err := DB.GetUserByEmail(ctx, email)
	if err != nil {
		return UserLoginResponse{}, fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return UserLoginResponse{}, fmt.Errorf("invalid password")
	}

	refreshTokens, err := DB.GetUserRefreshTokens(ctx, user.ID)
	if err != nil || len(refreshTokens) == 0 {
		return UserLoginResponse{}, fmt.Errorf("failed to get refresh tokens. Contact an admin to set a new refresh token")
	}

	return NewJWTFromRefreshToken(ctx, DB, refreshTokens[0].Token, jwtSecret)
}

func NewJWTFromRefreshToken(ctx context.Context, DB *db.Queries, refreshToken string, jwtSecret string) (UserLoginResponse, error) {
	user, err := DB.GetRefreshTokenUser(ctx, refreshToken)
	if err != nil {
		return UserLoginResponse{}, err
	}

	accessToken, err := CreateJWT(refreshToken, jwtSecret, time.Hour*1, DB, ctx)
	if err != nil {
		return UserLoginResponse{}, err
	}

	return UserLoginResponse{
		User: User{
			ID:      user.ID,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
		AccessToken:  accessToken,
		ExpiresAt:    time.Now().Add(time.Hour * 1),
		RefreshToken: refreshToken,
	}, nil
}
