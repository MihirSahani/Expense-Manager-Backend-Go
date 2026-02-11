package emiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/krakn/expense-management-backend-go/api/error"
	"github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/authenticator"
)

func Authenticate(authenticator authenticator.Authenticator, logger elogger.Logger, LOGGED_IN_USER_ID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getTokenFromHeader(r)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userId, err := authenticator.ValidateToken(token)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), LOGGED_IN_USER_ID, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTokenFromHeader(r *http.Request) (string, error) {
	// Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", eerror.MissingAuthenticationError
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", eerror.InvalidAuthenticationError
	}

	return authHeader[len(prefix):], nil
}