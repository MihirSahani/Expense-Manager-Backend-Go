package postgres_user

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/datastore"
)

func (p *PostgresUserStorage) DeleteUser(ctx context.Context, tx datastore.Database, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
