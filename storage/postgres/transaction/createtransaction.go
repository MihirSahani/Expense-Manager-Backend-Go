package postgres_transaction

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresTransactionStorage) CreateTransaction(ctx context.Context, tx *sql.Tx, t *entity.Transaction) (int64, error) {
	query := `
		INSERT INTO transactions (
			user_id, account_id, category_id, type, amount, payee, currency, 
			transaction_date, description, receipt_url, location
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, COALESCE(NULLIF($8::text, '')::DATE, NOW()), $9, $10, $11
		) RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query,
		t.UserID, t.AccountID, t.CategoryID, t.Type, t.Amount, t.Payee, t.Currency,
		t.TransactionDate, t.Description, t.ReceiptURL, t.Location,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
