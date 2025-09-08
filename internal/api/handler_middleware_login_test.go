package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CTSDM/motbwa-be/internal/auth"
	"github.com/CTSDM/motbwa-be/internal/database"
)

func TestHandlerMiddlewareLogin(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name           string
		setupUser      bool
		expiredToken   bool
		expectedStatus int
	}{
		{
			name:           "happy path - valid token with valid username",
			setupUser:      true,
			expiredToken:   false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "expired token",
			setupUser:      true,
			expiredToken:   true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "token signed with antoher secret",
			setupUser:      true,
			expiredToken:   true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing header token",
			setupUser:      false,
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

			var tokenString string
			if tc.setupUser {
				userData := database.CreateUserParams{
					Username:       "testuser",
					HashedPassword: "hashedpassword",
					CreatedAt:      time.Now().UTC(),
					UpdatedAt:      time.Now().UTC(),
				}
				// Create user and tokens
				user := createUser(t, &cfg, ctx, userData)
				token, err := auth.MakeJWT(user.ID, cfg.TokenSecret, cfg.TokenExpiration)
				if err != nil {
					t.Fatalf("failed to create JWT: %s", err)
				}
				if err != nil {
					t.Fatalf("unable to create tokens: %s", err)
				}
				tokenString = token
			}

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte{}))
			req.Header.Set("Content-Type", "application/json")
			req.Header["Auth"] = []string{fmt.Sprintf("Bearer %s", tokenString)}
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
