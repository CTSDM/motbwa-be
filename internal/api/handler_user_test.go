package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CTSDM/motbwa-be/internal/auth"
	"github.com/CTSDM/motbwa-be/internal/database"
)

func TestHandlerUser(t *testing.T) {
	testCases := []struct {
		name           string
		setupUser      bool
		emptyBody      bool
		expectedStatus int
	}{
		{
			name:           "happy path",
			setupUser:      false,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "user already exists",
			setupUser:      true,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "empty request body",
			setupUser:      false,
			emptyBody:      true,
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

			var reqData any
			var userData database.CreateUserParams
			username, password := "user", "test"
			reqData = parametersLogin{
				Username: username,
				Password: password,
			}

			if tc.setupUser {
				hashedPassword, err := auth.HashPassword(password)
				if err != nil {
					t.Fatalf("got unexpected error while hashing the password: %s", err)
				}
				userData = database.CreateUserParams{
					Username:       username,
					HashedPassword: hashedPassword,
					CreatedAt:      time.Now().UTC(),
					UpdatedAt:      time.Now().UTC(),
				}
				if _, err = cfg.DB.CreateUser(context.Background(), userData); err != nil {
					t.Fatalf("got unexpected error while creating the user in the database: %s", err)
				}

			}

			if tc.emptyBody {
				reqData = struct{}{}
			}

			reqBody, err := json.Marshal(reqData)
			if err != nil {
				t.Fatalf("something went wrong while marshaling the structure: %s", err)
			}
			req := httptest.NewRequestWithContext(context.Background(), "POST", "/test", bytes.NewReader(reqBody))
			req.Header.Set("Content Type", "application/json")
			rr := httptest.NewRecorder()

			cfg.HandlerCreateUser(rr, req)
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected %v status, got %v", tc.expectedStatus, rr.Code)
			}
		})
	}
}
