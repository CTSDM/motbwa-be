package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/CTSDM/motbwa-be/internal/auth"
	"github.com/google/uuid"
)

type userKey int

const (
	_ userKey = iota
	key
)

func ContextWithUser(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, key, userID)
}

func UserFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(key).(uuid.UUID)
	return userID, ok
}

func (cfg *CfgAPI) HandlerMiddlewareLogin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetHeaderValueTokenAPI(r.Header, "Auth")
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Token was not found in the header.", err)
			return
		}
		refreshTokenString, err := auth.GetHeaderValueTokenAPI(r.Header, "X-Refresh-Token")
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Token was not found in the header.", err)
			return
		}

		userID, err := auth.ValidateJWT(tokenString, cfg.TokenSecret)
		if err != nil {
			log.Println(refreshTokenString)
			refreshToken, err := cfg.DB.GetRefreshToken(context.Background(), refreshTokenString)
			if err == sql.ErrNoRows {
				RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
				return
			}
			if refreshToken.ExpiresAt.Before(time.Now()) {
				RespondWithError(w, http.StatusUnauthorized, "refresh token expired", nil)
				return
			}
			userID = refreshToken.UserID
			if time.Until(refreshToken.ExpiresAt) < time.Second*3600*6 {
				tokenString, refreshTokenString, err = cfg.createTokens(context.Background(), userID)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, "could not refresh the token", err)
					return
				}
			}
		}

		// add value to context
		ctx := r.Context()
		ctx = ContextWithUser(ctx, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
