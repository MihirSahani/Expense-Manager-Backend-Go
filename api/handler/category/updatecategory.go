package category

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

func UpdateCategory(logger elogger.Logger, storage *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get category id from URL
		categoryId, err := strconv.ParseInt(chi.URLParam(r, "categoryid"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Read category id from URL")
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

		// Validate payload
		if err := validate.Validate.Struct(payload); err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Payload validated")

		// Read from DB
		data, err := storage.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			category, err := storage.Category.GetCategoryByID(ctx, tx, categoryId, ctx.Value(LOGGED_IN_USER).(int64))
			if err != nil {
				return nil, err
			}
			if payload.Name != nil {
				category.Name = *payload.Name
			}
			if payload.Desc != nil {
				category.Desc = *payload.Desc
			}
			if payload.Type != nil {
				category.Type = *payload.Type
			}
			if payload.Color != nil {
				category.Color = *payload.Color
			}
			err = storage.Category.UpdateCategory(ctx, tx, category, ctx.Value(LOGGED_IN_USER).(int64))
			if err != nil {
				return nil, err
			}
			return category, nil
		})
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Read from database")
		// Return category
		ehandler.WriteJSON(w, http.StatusOK, data.(*entity.Category))
	})
}
