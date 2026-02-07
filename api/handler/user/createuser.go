package ehandleruser

import (
	"context"
	"database/sql"
	"net/http"

	elogger "github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/validate"
	"github.com/krakn/expense-management-backend-go/storage"
	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func CreateUser(logger elogger.Logger, s *storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the payload
		var payload struct {
			Email     string `json:"email" validate:"required,email"`
			FirstName string `json:"first_name" validate:"required"`
			LastName  string `json:"last_name"`
			Password  string `json:"password" validate:"required,min=8"`
		}
		err := ReadJSON(r, &payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Warn(err.Error())
			return
		}
		logger.Debug("Fetched from user details from payload")

		// validate the payload
		err = validate.Validate.Struct(payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Warn(err.Error())
			return
		}
		logger.Debug("Validation of payload passed")


		// hash the password
		hashedPassword, err := hashPassword([]byte(payload.Password))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Warn(err.Error())
			return
		}
		payload.Password = string(hashedPassword)
		logger.Debug("Password hashed")

		// create the user
		data, err := s.WithTransaction(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			userId, err := s.User.CreateUser(ctx, tx, entity.User{
				Email:     payload.Email,
				FirstName: payload.FirstName,
				LastName:  payload.LastName,
				Password:  payload.Password,
			})
			if err != nil {
				return nil, err
			}
			return userId, nil
		})
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Debug("User created")


		// return the user
		WriteJSON(w, http.StatusCreated, map[string]int64{
			"id": data.(int64),
		})
	})
}
