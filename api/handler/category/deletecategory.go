package category

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/storage"
)

func DeleteCategory(logger elogger.Logger, storage *storage.Storage) http.HandlerFunc {
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
		_, err = storage.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			err := storage.Category.DeleteCategory(ctx, tx, categoryId)
			if err != nil {
				return nil, err
			}
			return nil, nil
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