package postgres_transaction

import (
	"context"
	"database/sql"
)

func (s *PostgresTransactionStorage) DeleteTransaction(ctx context.Context, tx *sql.Tx, id int64, userId int64) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`
	_, err := tx.ExecContext(ctx, query, id, userId)
	return err
}
