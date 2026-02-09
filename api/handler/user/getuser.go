package ehandleruser

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"github.com/krakn/expense-management-backend-go/storage/entity"
	"go.uber.org/zap"
)

func GetuserByID(logger elogger.Logger, s *storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the id from the request
		userId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if userId == 0 {
			logger.Warn("Client did not provide id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Fetched user id from URL", zap.Int64("id", userId))

		// read from database
		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx datastore.Database) (any, error) {
			return s.User.GetUserByID(ctx, tx, userId)
		})

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user := data.(entity.User)
		logger.Debug("Fetched user from database", zap.Any("user", user))

		// return user info
		ehandler.WriteJSON(w, http.StatusOK, user)
	})
}
