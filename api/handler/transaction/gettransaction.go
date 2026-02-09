package transaction

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/datastore"
)

func GetAllTransactions(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx datastore.Database) (any, error) {
			return s.Transaction.GetAllTransactions(ctx, tx, ctx.Value(LOGGED_IN_USER).(int64))
		})
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ehandler.WriteJSON(w, http.StatusOK, data)
	})
}

func GetTransactionByID(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx datastore.Database) (any, error) {
			return s.Transaction.GetTransactionByID(ctx, tx, id, ctx.Value(LOGGED_IN_USER).(int64))
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
