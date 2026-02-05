package app

import (
	"net/http"
	"strings"
	"context"
	"github.com/krakn/expense-management-backend-go/api/application-server/error"
	"go.uber.org/zap"
)

func (a *application) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromHeader(r)
		if err != nil {
			a.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		userId, err := a.authenticator.ValidateToken(token)
		if err != nil {
			a.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), LOGGED_IN_USER_ID, userId)

		a.logger.Info("Authenticated User", zap.Int64("user_id", userId))
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getTokenFromHeader(r *http.Request) (string, error) {
	// Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.MissingAuthenticationError
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.InvalidAuthenticationError
	}

	return authHeader[len(prefix):], nil
}