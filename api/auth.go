package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hydrocode-de/gorun/internal/auth"
	"github.com/hydrocode-de/gorun/internal/db"
	"github.com/spf13/viper"
)

func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	// the refresh token is sent as a JSON body
	var refreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body. A refresh token is required.")
	}

	DB := viper.Get("db").(*db.Queries)
	secret := viper.GetString("secret")

	response, err := auth.NewJWTFromRefreshToken(r.Context(), DB, refreshToken.RefreshToken, secret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid refresh token: %v", err))
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse request body: %v", err))
		return
	}

	DB := viper.Get("db").(*db.Queries)
	secret := viper.GetString("secret")
	response, err := auth.LoginUser(r.Context(), DB, loginRequest.Email, loginRequest.Password, secret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Login attempt failed: %v", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, response)
}
