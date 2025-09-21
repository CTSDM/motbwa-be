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

func TestHandlerLogin(t *testing.T) {
	testCases := []struct {
		name           string
		setupUser      bool
		invalidPayload bool
		emptyPayload   bool
		expectedStatus int
	}{
		{
			name:           "happy_path",
			setupUser:      true,
			invalidPayload: false,
			emptyPayload:   false,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "not existing user",
			setupUser:      false,
			invalidPayload: false,
			emptyPayload:   false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "payload with wrong payload",
			setupUser:      false,
			invalidPayload: true,
			emptyPayload:   false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty payload",
			setupUser:      false,
			invalidPayload: false,
			emptyPayload:   true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the database
			t.Cleanup(cleanup)

			cfg := CfgAPI{
				DB:                     database.New(db),
				TokenSecret:            "test-secret",
				TokenExpiration:        time.Second * 60,
				RefreshTokenExpiration: time.Hour * 24,
			}

			var reqData any
			username, password := "user", "test"
			var user database.User

			reqData = parametersLogin{
				Username: username,
				Password: password,
			}

			if tc.setupUser {
				hashedPassword, err := auth.HashPassword(password)
				if err != nil {
					t.Fatalf("couldn't hash the password :%s", err)
				}
				user, err = cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
					Username:       username,
					HashedPassword: hashedPassword,
					CreatedAt:      time.Now().UTC(),
					UpdatedAt:      time.Now().UTC(),
				})
				if err != nil {
					t.Fatalf("something went wrong while creating the user in the database: %s", err)
				}
			}

			if tc.invalidPayload {
				reqData = struct{ Username int }{Username: 10}
			}
			if tc.emptyPayload {
				reqData = struct{}{}
			}

			reqBody, err := json.Marshal(reqData)
			if err != nil {
				t.Fatalf("couldn't marshal the request body: %s", err)
			}
			// Setup request and response recorder
			req := httptest.NewRequestWithContext(context.Background(), "POST", "/tets", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Call the function to test
			cfg.HandlerLogin(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, rr.Code)
				return
			}
			if tc.expectedStatus == http.StatusCreated {
				var res responseValsLogin
				decoder := json.NewDecoder(rr.Body)
				if err := decoder.Decode(&res); err != nil {
					t.Fatalf("couldn't decode the response payload: %s", err)
				}

				if user.ID != res.ID {
					t.Errorf("got userID %s, want userID %s", res.ID, user.ID)
				}
				if user.Username != res.Username {
					t.Errorf("got username %s, want username %s", res.Username, user.Username)
				}
			}
		})
	}
}
