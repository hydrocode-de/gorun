package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hydrocode-de/gorun/config"
	"github.com/hydrocode-de/gorun/internal/auth"
)

func HandleRefreshToken(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	// the refresh token is sent as a JSON body
	var refreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body. A refresh token is required.")
	}

	response, err := auth.NewJWTFromRefreshToken(r.Context(), c.GetDB(), refreshToken.RefreshToken, c.Secret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid refresh token: %v", err))
	}

	ResondWithJSON(w, http.StatusOK, response)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, c *config.APIConfig) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse request body: %v", err))
		return
	}

	response, err := auth.LoginUser(r.Context(), c.GetDB(), loginRequest.Email, loginRequest.Password, c.Secret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Login attempt failed: %v", err))
		return
	}

	ResondWithJSON(w, http.StatusOK, response)
}
