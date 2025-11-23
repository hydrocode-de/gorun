package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/hydrocode-de/gorun/internal/helper"
	"github.com/spf13/viper"
)

type AdminCredentials struct {
	UserID       string    `json:"userId"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refreshToken"`
	AccessToken  string    `json:"accessToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func createNewAdminUser(ctx context.Context) (UserLoginResponse, error) {
	DB := viper.Get("db").(*db.Queries)
	JWTSecret := viper.GetString("secret")

	// Generate a random password for the admin user
	password := helper.GetRandomString(16)

	// Create the admin user
	response, err := CreateUser(
		ctx,
		DB,
		"admin@gorun.local",
		password,
		true, // isAdmin
		JWTSecret,
	)

	if err != nil {
		return UserLoginResponse{}, err
	}

	// Print the admin password to the console
	fmt.Printf("Admin user created with password: %s\n", password)
	fmt.Println("Please save this password securely. It will not be shown again.")

	return response, nil
}

func CreateAdminCredentials(ctx context.Context) (*AdminCredentials, error) {
	basePath := viper.GetString("path")
	DB := viper.Get("db").(*db.Queries)

	// Check if admin credentials file exists
	credentialsPath := filepath.Join(basePath, "admin_credentials.json")

	// Try to load existing credentials
	if _, err := os.Stat(credentialsPath); err == nil {
		// File exists, try to load it
		credentials, err := GetAdminCredentials(ctx)
		if err == nil && credentials != nil {
			// Check if the credentials are still valid
			if time.Now().Before(credentials.ExpiresAt) {
				return credentials, nil
			}
		}
	}

	// Check if admin user exists in the database
	var adminResponse UserLoginResponse
	adminUser, err := DB.GetUserByEmail(ctx, "admin@gorun.local")

	if err != nil {
		// Admin user doesn't exist, create it
		adminResponse, err = createNewAdminUser(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create admin user: %w", err)
		}
	} else {
		// Admin user exists, create a new refresh token for it
		JWTSecret := viper.GetString("secret")

		// Create a new refresh token
		refreshToken := helper.GetRandomString(32)
		_, err = DB.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			UserID:    adminUser.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create refresh token for admin user: %w", err)
		}

		// Get login response using the new refresh token
		adminResponse, err = NewJWTFromRefreshToken(ctx, DB, refreshToken, JWTSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to create admin credentials: %w", err)
		}
	}

	// Create admin credentials
	credentials := &AdminCredentials{
		UserID:       adminResponse.User.ID,
		Email:        adminResponse.User.Email,
		RefreshToken: adminResponse.RefreshToken,
		AccessToken:  adminResponse.AccessToken,
		ExpiresAt:    adminResponse.ExpiresAt,
	}

	// Save credentials to file
	credentialsJSON, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal admin credentials: %w", err)
	}

	err = os.WriteFile(credentialsPath, credentialsJSON, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to save admin credentials: %w", err)
	}

	return credentials, nil
}

func GetAdminCredentials(ctx context.Context) (*AdminCredentials, error) {
	basePath := viper.GetString("path")
	credentialsPath := filepath.Join(basePath, "admin_credentials.json")

	// Check if file exists
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("admin credentials file not found")
	}

	// Read the file
	credentialsJSON, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read admin credentials: %w", err)
	}

	// Parse the JSON
	var credentials AdminCredentials
	err = json.Unmarshal(credentialsJSON, &credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to parse admin credentials: %w", err)
	}

	// Check if credentials are expired
	if time.Now().After(credentials.ExpiresAt) {
		DB := viper.Get("db").(*db.Queries)
		JWTSecret := viper.GetString("secret")

		// Try to refresh the token
		response, err := NewJWTFromRefreshToken(ctx, DB, credentials.RefreshToken, JWTSecret)
		if err != nil {
			return nil, fmt.Errorf("admin credentials expired and refresh failed: %w", err)
		}

		// Update credentials
		credentials.AccessToken = response.AccessToken
		credentials.ExpiresAt = response.ExpiresAt

		// Save updated credentials
		credentialsJSON, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal updated admin credentials: %w", err)
		}

		err = os.WriteFile(credentialsPath, credentialsJSON, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to save updated admin credentials: %w", err)
		}
	}

	return &credentials, nil
}
