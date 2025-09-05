package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CTSDM/motbwa-be/internal/database"
	"github.com/google/uuid"
)

func TestHandlerMiddlewareLogin(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name                string
		setupUser           bool
		expiredToken        bool
		expiredRefreshToken bool
		invalidUserID       bool
		wrongPayload        bool
		expectedStatus      int
	}{
		{
			name:                "happy path - valid token with valid username",
			setupUser:           true,
			expiredToken:        false,
			expiredRefreshToken: false,
			invalidUserID:       false,
			expectedStatus:      http.StatusOK,
		},
		{
			name:                "expired token with valid refresh",
			setupUser:           true,
			expiredToken:        true,
			expiredRefreshToken: false,
			invalidUserID:       false,
			expectedStatus:      http.StatusOK,
		},
		{
			name:                "expired token with expired refresh",
			setupUser:           true,
			expiredToken:        true,
			expiredRefreshToken: true,
			invalidUserID:       false,
			expectedStatus:      http.StatusUnauthorized,
		},
		{
			name:           "invalid user ID",
			setupUser:      true,
			expiredToken:   false,
			invalidUserID:  true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "no user setup",
			setupUser:      false,
			wrongPayload:   false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "wrong payload",
			setupUser:      false,
			wrongPayload:   true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean database before each test
			t.Cleanup(cleanup)

			cfg := CfgAPI{
				DB:                     database.New(db),
				TokenSecret:            "test-secret",
				TokenExpiration:        time.Second * 60,
				RefreshTokenExpiration: time.Hour * 24,
			}

			if tc.expiredToken {
				cfg.TokenExpiration = 0
			}
			if tc.expiredRefreshToken {
				cfg.RefreshTokenExpiration = 0
			}

			var reqData any

			if tc.setupUser {
				userData := database.CreateUserParams{
					Username:       "testuser",
					HashedPassword: "hashedpassword",
					CreatedAt:      time.Now().UTC(),
					UpdatedAt:      time.Now().UTC(),
				}
				// Create user and tokens
				user := createUser(t, &cfg, ctx, userData)
				tokenString, refreshTokenString, err := cfg.createTokens(ctx, user.ID)
				if err != nil {
					t.Fatalf("unable to create tokens: %s", err)
				}

				userID := user.ID
				if tc.invalidUserID {
					userID = uuid.New()
				}

				reqData = authBody{
					UserID:             userID,
					Username:           user.Username,
					TokenString:        tokenString,
					RefreshTokenString: refreshTokenString,
				}

			} else if tc.wrongPayload {
				reqData = struct {
					UserID string
				}{
					UserID: "INVALID",
				}
			} else {
				reqData = authBody{
					UserID:             uuid.New(),
					Username:           "none",
					TokenString:        "none",
					RefreshTokenString: "none",
				}
			}

			// Marshal the request data to JSON
			reqBody, err := json.Marshal(reqData)
			if err != nil {
				t.Fatalf("failed to marshal request body: %s", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			// Create response recorder
			rr := httptest.NewRecorder()

			// Create a dummy next
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			handler := cfg.HandlerMiddlewareLogin(nextHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func createUser(t *testing.T, cfg *CfgAPI, ctx context.Context, userData database.CreateUserParams) database.User {
	user, err := cfg.DB.CreateUser(ctx, userData)
	if err != nil {
		t.Helper()
		t.Fatalf("couldn't create the user: %s", err)
	}
	return user
}
