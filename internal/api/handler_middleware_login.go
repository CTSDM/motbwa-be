package api

import (
	"context"
	"net/http"

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
		// get the token from the headers
		tokenString, err := auth.GetHeaderValueTokenAPI(r.Header, "Auth")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Token was not found in the header.", err)
			return
		}
		userID, err := auth.ValidateJWT(tokenString, cfg.TokenSecret)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
			return
		}

		// add value to context
		ctx := r.Context()
		ctx = ContextWithUser(ctx, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
