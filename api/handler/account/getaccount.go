package account

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/storage"
)

func GetAllAccounts(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(LOGGED_IN_USER).(int64)

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			return s.Account.GetAllAccounts(ctx, tx, userID)
		})

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ehandler.WriteJSON(w, http.StatusOK, data)
	})
}

func GetAccountByID(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId := r.Context().Value(LOGGED_IN_USER).(int64)

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			acc, err := s.Account.GetAccountByID(ctx, tx, id, userId)
			if err != nil {
				return nil, err
			}
			if acc.UserID != userId {
				return nil, sql.ErrNoRows
			}
			return acc, nil
		})

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ehandler.WriteJSON(w, http.StatusOK, data)
	})
}
