package transaction

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"go.uber.org/zap"
)

func GetAllTransactions(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Page     *int64 `json:"page" validate:"omitempty,min=1"`
			PageSize *int64 `json:"page_size" validate:"omitempty,min=1,max=100"`
		}
		if err := ehandler.ReadJSON(r, &payload); err != nil && err != io.EOF {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Read from payload successfully")

		if err := validate.Validate.Struct(payload); err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Successfully validated payload")

		totalTransactions, err := s.Transaction.GetTransactionCount(r.Context(), s.Connection.GetDb(), r.Context().Value(LOGGED_IN_USER).(int64))
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully fetched total transactions")

		page := int64(1)
		if payload.Page == nil {
			page = 1
		} else {
			page = *payload.Page
		}
		pageSize := int64(10)
		if payload.PageSize == nil {
			pageSize = 10
		} else {
			pageSize = *payload.PageSize
		}

		data, err := s.Transaction.GetTransactionsPaginated(r.Context(), s.Connection.GetDb(), r.Context().Value(LOGGED_IN_USER).(int64), page, pageSize)

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully fetched transactions")


		totalPages := (totalTransactions + pageSize - 1) / pageSize
		ehandler.WriteJSON(w, http.StatusOK, map[string]any{
			"transactions":       data,
			"total_pages":        totalPages,
			"total_transactions": totalTransactions,
		})
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

func GetTransactionsByMonth(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Month int64 `json:"month" validate:"required,min=1,max=12"`
			Year  int64 `json:"year" validate:"required,min=1900,max=2026"`
		}
		if err := ehandler.ReadJSON(r, &payload); err != nil {
			logger.Warn(err.Error(), zap.Any("payload", payload))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Read from payload successfully")

		if err := validate.Validate.Struct(payload); err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Successfully validated payload")

		data, err := s.Transaction.GetTransactionsByMonth(r.Context(), s.Connection.GetDb(), r.Context().Value(LOGGED_IN_USER).(int64), payload.Month, payload.Year)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully fetched transactions")

		ehandler.WriteJSON(w, http.StatusOK, data)
	})
}
