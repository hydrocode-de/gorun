package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/helper"
)

type AdminCredentials struct {
	UserID       string    `json:"userId"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refreshToken"`
	AccessToken  string    `json:"accessToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func createNewAdminUser(c *APIConfig) (auth.UserLoginResponse, error) {
	// Generate a random password for the admin user
	password := helper.GetRandomString(16)

	// Create the admin user
	response, err := auth.CreateUser(
		context.Background(),
		c.GetDB(),
		"admin@gorun.local",
		password,
		true, // isAdmin
		c.Secret,
	)

	if err != nil {
		return auth.UserLoginResponse{}, err
	}

	// Print the admin password to the console
	fmt.Printf("Admin user created with password: %s\n", password)
	fmt.Println("Please save this password securely. It will not be shown again.")

	return response, nil
}

func (c *APIConfig) CreateAdminCredentials() (*AdminCredentials, error) {
	// Check if admin credentials file exists
	credentialsPath := filepath.Join(c.GorunBasePath, "admin_credentials.json")

	// Try to load existing credentials
	if _, err := os.Stat(credentialsPath); err == nil {
		// File exists, try to load it
		credentials, err := c.GetAdminCredentials()
		if err == nil && credentials != nil {
			// Check if the credentials are still valid
			if time.Now().Before(credentials.ExpiresAt) {
				return credentials, nil
			}
		}
	}

	// Check if admin user exists in the database
	ctx := context.Background()
	var adminResponse auth.UserLoginResponse
	_, err := c.GetDB().GetUserByEmail(ctx, "admin@gorun.local")

	if err != nil {
		// Admin user doesn't exist, create it
		adminResponse, err = createNewAdminUser(c)
		if err != nil {
			return nil, fmt.Errorf("failed to create admin user: %w", err)
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

func (c *APIConfig) GetAdminCredentials() (*AdminCredentials, error) {
	credentialsPath := filepath.Join(c.GorunBasePath, "admin_credentials.json")

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
		// Try to refresh the token
		ctx := context.Background()
		response, err := auth.NewJWTFromRefreshToken(ctx, c.GetDB(), credentials.RefreshToken, c.Secret)
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
