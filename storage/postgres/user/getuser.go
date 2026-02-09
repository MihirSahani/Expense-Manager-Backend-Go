package postgres_user

import (
	"context"

	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (p *PostgresUserStorage) GetUserByEmail(ctx context.Context, tx datastore.Database, email string) (entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user entity.User
	err := tx.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (p *PostgresUserStorage) GetUserByID(ctx context.Context, tx datastore.Database, id int64) (entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entity.User
	err := tx.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
