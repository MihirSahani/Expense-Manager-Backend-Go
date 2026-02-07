package ehandleruser

import (
	"context"
	"database/sql"
	"net/http"

	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/authenticator"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/entity"
	"golang.org/x/crypto/bcrypt"
)

func Login(logger elogger.Logger, s *storage.Storage, a authenticator.Authenticator) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := ReadJSON(r, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Debug("Fetched from user details from payload")

		if err = validate.Validate.Struct(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Debug("Validation of payload passed")

		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			user, err := s.User.GetUserByEmail(ctx, tx, payload.Email)
			if err != nil {
				return nil, err
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
			if err != nil {
				return nil, err
			}

			return user, nil
		})
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		logger.Debug("User authenticated")
		user := data.(entity.User)

		token, err := a.GenerateToken(user.Id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Debug("Token generated")

		ehandler.WriteJSON(w, http.StatusOK, map[string]string{
			"token": token,
		})
	})
}
