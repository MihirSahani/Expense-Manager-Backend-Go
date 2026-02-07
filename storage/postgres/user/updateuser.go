package postgres_user

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (p *PostgresUserStorage) UpdateUser(ctx context.Context, tx *sql.Tx, user entity.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := tx.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Id)
	if err != nil {
		return err
	}

	return nil
}