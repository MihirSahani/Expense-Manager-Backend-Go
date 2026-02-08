package transaction

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/entity"
	"go.uber.org/zap"
)

func UpdateTransaction(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var payload struct {
			AccountID       *int64   `json:"account_id"`
			CategoryID      *int64   `json:"category_id"`
			Type            *string  `json:"type" validate:"omitempty,oneof=income expense transfer"`
			Amount          *float64 `json:"amount" validate:"omitempty,gt=0"`
			Payee           *string  `json:"payee"`
			Currency        *string  `json:"currency" validate:"omitempty,len=3"`
			TransactionDate *string  `json:"transaction_date"`
			Description     *string  `json:"description"`
			ReceiptURL      *string  `json:"receipt_url"`
			Location        *string  `json:"location"`
		}

		if err := ehandler.ReadJSON(r, &payload); err != nil {
			logger.Warn(err.Error(), zap.Any("payload", payload))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := validate.Validate.Struct(payload); err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			t, err := s.Transaction.GetTransactionByID(ctx, tx, id, ctx.Value(LOGGED_IN_USER).(int64))
			if err != nil {
				return nil, err
			}

			if payload.AccountID != nil {
				t.AccountID = *payload.AccountID
			}
			if payload.CategoryID != nil {
				t.CategoryID = *payload.CategoryID
			}
			if payload.Type != nil {
				t.Type = *payload.Type
			}
			if payload.Amount != nil {
				t.Amount = *payload.Amount
			}
			if payload.Payee != nil {
				t.Payee = *payload.Payee
			}
			if payload.Currency != nil {
				t.Currency = *payload.Currency
			}
			if payload.TransactionDate != nil {
				t.TransactionDate = *payload.TransactionDate
			}
			if payload.Description != nil {
				t.Description = *payload.Description
			}
			if payload.ReceiptURL != nil {
				t.ReceiptURL = *payload.ReceiptURL
			}
			if payload.Location != nil {
				t.Location = *payload.Location
			}

			return t, s.Transaction.UpdateTransaction(ctx, tx, t)
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

		ehandler.WriteJSON(w, http.StatusOK, data.(*entity.Transaction))
	})
}
