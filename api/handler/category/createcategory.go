package category

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

func CreateCategory(logger elogger.Logger, storage *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read payload
		var payload struct {
			Name  *string `json:"name" validate:"required,min=3,max=100"`
			Desc  *string `json:"desc" validate:"required,min=3,max=100"`
			Type  *string `json:"type" validate:"required,min=3,max=100,oneof=income expense"`
			Color *string `json:"color"`
		}
		if err := ehandler.ReadJSON(r, &payload); err != nil {
			logger.Warn(err.Error(), zap.Any("payload", payload))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Fetched payload successfully")

		// Validate the paylaod
		if err := validate.Validate.Struct(payload); err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Payload validated")

		// Write to DB
		data, err := storage.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			id, err := storage.Category.CreateCategory(ctx, tx, &entity.Category{
				Name:   *payload.Name,
				Desc:   *payload.Desc,
				Type:   *payload.Type,
				Color:  *payload.Color,
				UserID: ctx.Value(LOGGED_IN_USER).(int64),
			})
			if err != nil {
				return nil, err
			}
			return id, nil
		})
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Wrote to database")
		// Return id
		ehandler.WriteJSON(w, http.StatusCreated, map[string]int64{"id": data.(int64)})
	})
}
