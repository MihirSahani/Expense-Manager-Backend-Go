package ehandleruser

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"go.uber.org/zap"
)

func DeleteUser(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the id from the request
		userId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Fetched user id from URL", zap.Int64("id", userId))

		// match it with the logged in user
		if userId != r.Context().Value(LOGGED_IN_USER).(int64) {
			logger.Warn("User is not authorized to delete this user")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Debug("User is authorized to delete this user")

		// delete the user
		_, err = s.WithTransaction(r.Context(), func(ctx context.Context, tx datastore.Database) (any, error) {
			return nil, s.User.DeleteUser(ctx, tx, userId)
		})
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("User deleted")
		w.WriteHeader(http.StatusOK)
	})
}
