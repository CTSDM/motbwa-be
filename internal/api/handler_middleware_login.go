package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CTSDM/motbwa-be/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authBody struct {
	UserID             uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	TokenString        string    `json:"tokenString"`
	RefreshTokenString string    `json:"refreshTokenString"`
}

func (cfg *CfgAPI) HandlerMiddlewareLogin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params authBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, err := auth.ValidateJWT(params.TokenString, cfg.TokenSecret)
		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshTokenDB, err2 := cfg.DB.GetRefreshToken(r.Context(), params.UserID)
			if err2 == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err2 != nil {
				RespondWithError(w, http.StatusInternalServerError, "failed to retrieve refresh token", err2)
				return
			}
			// we check if the refresh token has expired
			timeNow := time.Now()
			if timeNow.After(refreshTokenDB.ExpiresAt) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else if err != nil || params.UserID != userID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
