package transaction

import (
	"context"
	"database/sql"
	"net/http"

	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/entity"
	"go.uber.org/zap"
)

func CreateTransaction(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			AccountID       *int64   `json:"account_id" validate:"required"`
			CategoryID      *int64   `json:"category_id" validate:"required"`
			Type            *string  `json:"type" validate:"required,oneof=income expense transfer"`
			Amount          *float64 `json:"amount" validate:"required,gt=0"`
			Payee           *string  `json:"payee" validate:"required"`
			Currency        *string  `json:"currency" validate:"required,len=3"`
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
		logger.Debug("Fetched payload from request")

		if err := validate.Validate.Struct(payload); err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Successfully validated payload")

		description := ""
		if payload.Description != nil {
			description = *payload.Description
		}
		receiptURL := ""
		if payload.ReceiptURL != nil {
			receiptURL = *payload.ReceiptURL
		}
		location := ""
		if payload.Location != nil {
			location = *payload.Location
		}
		transactionDate := ""
		if payload.TransactionDate != nil {
			transactionDate = *payload.TransactionDate
		}

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			return s.Transaction.CreateTransaction(ctx, tx, &entity.Transaction{
				UserID:          ctx.Value(LOGGED_IN_USER).(int64),
				AccountID:       *payload.AccountID,
				CategoryID:      *payload.CategoryID,
				Type:            *payload.Type,
				Amount:          *payload.Amount,
				Payee:           *payload.Payee,
				Currency:        *payload.Currency,
				TransactionDate: transactionDate,
				Description:     description,
				ReceiptURL:      receiptURL,
				Location:        location,
			})
		})

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully created transaction")

		ehandler.WriteJSON(w, http.StatusCreated, map[string]int64{"id": data.(int64)})
	})
}
