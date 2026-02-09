package category

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/datastore"
)

func DeleteCategory(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get Id from URL
		categoryId, err := strconv.ParseInt(chi.URLParam(r, "categoryid"), 10, 64)
		if err != nil {
			logger.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Debug("Read category Id from URL")

		// delete from Db
		_, err = s.WithTransaction(r.Context(), func(ctx context.Context, tx datastore.Database) (any, error) {
			return nil, s.Category.DeleteCategory(ctx, tx, categoryId, ctx.Value(LOGGED_IN_USER).(int64))
		})
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				logger.Warn(err.Error())
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				logger.Warn(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		logger.Debug("Delete from database")

		// return success
		w.WriteHeader(http.StatusOK)
	})
}
