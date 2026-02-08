package postgres_account

import (
	"context"
	"database/sql"
)

func (s *PostgresAccountStorage) DeleteAccount(ctx context.Context, tx *sql.Tx, id int64, userId int64) error {
	query := `
		DELETE FROM 
			accounts 
		WHERE 
			id = $1 AND user_id = $2
	`
	_, err := tx.ExecContext(ctx, query, id, userId)
	return err
}
