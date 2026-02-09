package category

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
)

func GetAllCategory(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read from DB
		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, db datastore.Database) (any, error) {
			return s.Category.GetAllCategories(ctx, db, ctx.Value(LOGGED_IN_USER).(int64))
		})
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Read from database")
		// Return categories
		ehandler.WriteJSON(w, http.StatusOK, data.(*[]entity.Category))
	})
}

func GetCategoryByID(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Id from URL
		categoryId, err := strconv.ParseInt(chi.URLParam(r, "categoryid"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Read category Id from URL")

		// Read from DB
		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, db datastore.Database) (any, error) {
			return s.Category.GetCategoryByID(ctx, db, categoryId, ctx.Value(LOGGED_IN_USER).(int64))
		})
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debug("Read from database")
		// Return category
		ehandler.WriteJSON(w, http.StatusOK, data.(*entity.Category))
	})
}
