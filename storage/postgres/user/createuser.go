package postgres_user

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (p *PostgresUserStorage) CreateUser(ctx context.Context, tx *sql.Tx, user entity.User) (int64, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}