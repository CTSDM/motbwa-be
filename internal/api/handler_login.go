package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CTSDM/motbwa-be/internal/auth"
	"github.com/CTSDM/motbwa-be/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type parameters struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseVals struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

var ErrUnauthorized = fmt.Errorf("unauthorized")

func (cfg *CfgAPI) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// Get username from the body
	params, err := cfg.parseLoginRequest(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error(), err)
	}

	// Authenticate user
	user, err := cfg.authenticateUser(r.Context(), *params)
	if err == ErrUnauthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}

	// Generate tokens
	tokenString, refreshTokenString, err := cfg.createTokens(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// All is good so let's return the jwt
	RespondWithJSON(w, http.StatusCreated, responseVals{
		ID:           user.ID,
		Username:     user.Username,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	})
}

func (cfg *CfgAPI) authenticateUser(ctx context.Context, params parameters) (*database.User, error) {
	// Get user from DB
	user, err := cfg.DB.GetUserByUsername(ctx, params.Username)
	if err == sql.ErrNoRows {
		return nil, ErrUnauthorized
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Check the password
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrUnauthorized
	} else if err != nil {
		return nil, fmt.Errorf("failed to check password: %w", err)
	}

	return &user, nil
}

func (cfg *CfgAPI) createTokens(ctx context.Context, userID uuid.UUID) (string, string, error) {
	// Create tokens
	tokenString, err := auth.MakeJWT(userID, cfg.TokenSecret, cfg.TokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed to create JWT: %w", err)
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	// When the user logs in, always give it a new refresh token
	refreshToken, err := cfg.DB.GetRefreshToken(ctx, userID)
	switch err {
	case sql.ErrNoRows:
		_, err = cfg.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
			Token:     refreshTokenString,
			UserID:    userID,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(cfg.RefreshTokenExpiration),
		})
	case nil:
		refreshToken, err = cfg.DB.UpdateRefreshToken(ctx, database.UpdateRefreshTokenParams{
			UserID:    userID,
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(cfg.RefreshTokenExpiration),
		})
	default:
		return "", "", fmt.Errorf("failed to deleted old refresh token: %w", err)
	}

	return tokenString, refreshToken.Token, err
}

func (cfg *CfgAPI) parseLoginRequest(r *http.Request) (*parameters, error) {
	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return nil, fmt.Errorf("failed to decode the body: %w", err)
	}
	return &params, nil
}
