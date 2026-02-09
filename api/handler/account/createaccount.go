package account

import (
	"context"
	"net/http"

	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"github.com/krakn/expense-management-backend-go/storage/entity"
	"go.uber.org/zap"
)

func CreateAccount(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Name              *string  `json:"name" validate:"required,min=3,max=100"`
			Type              *string  `json:"type" validate:"required,oneof=cash bank credit_card digital_wallet investment"`
			Currency          *string  `json:"currency" validate:"required,len=3"`
			CurrentBalance    *float64 `json:"current_balance" validate:"required"`
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

		// Defaults
		isIncluded := true
		if payload.IsIncludedInTotal != nil {
			isIncluded = *payload.IsIncludedInTotal
		}
		isActive := true
		if payload.IsActive != nil {
			isActive = *payload.IsActive
		}
		bankName := ""
		if payload.BankName != nil {
			bankName = *payload.BankName
		}
		accountNumber := ""
		if payload.AccountNumber != nil {
			accountNumber = *payload.AccountNumber
		}

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx datastore.Database) (any, error) {
			return s.Account.CreateAccount(ctx, tx, &entity.Account{
				Name:              *payload.Name,
				Type:              *payload.Type,
				Currency:          *payload.Currency,
				CurrentBalance:    *payload.CurrentBalance,
				BankName:          bankName,
				AccountNumber:     accountNumber,
				IsIncludedInTotal: isIncluded,
				UserID:            r.Context().Value(LOGGED_IN_USER).(int64),
				IsActive:          isActive,
			})
		})

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ehandler.WriteJSON(w, http.StatusCreated, map[string]int64{"id": data.(int64)})
	})
}
