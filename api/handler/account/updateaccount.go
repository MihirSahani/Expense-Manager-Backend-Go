package account

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
	"go.uber.org/zap"
)

func UpdateAccount(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var payload struct {
			Name              *string  `json:"name" validate:"omitempty,min=3,max=100"`
			Type              *string  `json:"type" validate:"omitempty,oneof=cash bank credit_card digital_wallet investment"`
			Currency          *string  `json:"currency" validate:"omitempty,len=3"`
			CurrentBalance    *float64 `json:"current_balance"`
			BankName          *string  `json:"bank_name"`
			AccountNumber     *string  `json:"account_number"`
			IsIncludedInTotal *bool    `json:"is_included_in_total"`
			IsActive          *bool    `json:"is_active"`
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

		userID := r.Context().Value(LOGGED_IN_USER).(int64)

		_, err = s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			existing, err := s.Account.GetAccountByID(ctx, tx, id, userID)
			if err != nil {
				return nil, err
			}
			if existing.UserID != userID {
				return nil, sql.ErrNoRows
			}

			if payload.Name != nil {
				existing.Name = *payload.Name
			}
			if payload.Type != nil {
				existing.Type = *payload.Type
			}
			if payload.Currency != nil {
				existing.Currency = *payload.Currency
			}
			if payload.CurrentBalance != nil {
				existing.CurrentBalance = *payload.CurrentBalance
			}
			if payload.BankName != nil {
				existing.BankName = *payload.BankName
			}
			if payload.AccountNumber != nil {
				existing.AccountNumber = *payload.AccountNumber
			}
			if payload.IsIncludedInTotal != nil {
				existing.IsIncludedInTotal = *payload.IsIncludedInTotal
			}
			if payload.IsActive != nil {
				existing.IsActive = *payload.IsActive
			}

			return nil, s.Account.UpdateAccount(ctx, tx, existing, userID)
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

		w.WriteHeader(http.StatusOK)
	})
}
