package ehandler

import (
	"context"
	"database/sql"
	"net/http"

	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func UpdateUser(logger elogger.Logger, s *storage.Storage, LOGGED_IN_USER string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read the payload
		var payload struct {
			Email     *string `json:"email"`
			FirstName *string `json:"first_name"`
			LastName  *string `json:"last_name"`
			Password  *string `json:"password"`
		}
		err := ReadJSON(r, &payload)
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Debug("Fetched from user details from payload")
		
		// validate the payload
		err = validate.Validate.Struct(payload)
		if err != nil {
			logger.Warn(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Debug("Validation of payload passed")

		// hash the password
		var hashedPassword []byte
		if payload.Password != nil {
			hashedPassword, err = hashPassword([]byte(*payload.Password))
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			*payload.Password = string(hashedPassword)
		}
		logger.Debug("Password hashed")

		// update the user
		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, x *sql.Tx) (any, error) {
			user, err := s.User.GetUserByID(ctx, x, ctx.Value(LOGGED_IN_USER).(int64))
			if err != nil {
				return nil, err
			}
			if payload.Email != nil {
				user.Email = *payload.Email
			}
			if payload.FirstName != nil {
				user.FirstName = *payload.FirstName
			}
			if payload.LastName != nil {
				user.LastName = *payload.LastName
			}
			if payload.Password != nil {
				user.Password = *payload.Password
			}

			err = s.User.UpdateUser(ctx, x, user)
			if err != nil {
				return nil, err
			}
			return user, nil
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Debug("User updated")
		
		// return the user
		WriteJSON(w, http.StatusOK, data.(entity.User))
	})
}