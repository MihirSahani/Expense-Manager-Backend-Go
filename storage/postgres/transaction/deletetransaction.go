package postgres_transaction

import (
	"context"

	"github.com/krakn/expense-management-backend-go/storage/datastore"
)

func (s *PostgresTransactionStorage) DeleteTransaction(ctx context.Context, tx datastore.Database, id int64, userId int64) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`
	_, err := tx.ExecContext(ctx, query, id, userId)
	return err
}
