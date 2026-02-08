package postgres_transaction

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresTransactionStorage) UpdateTransaction(ctx context.Context, tx *sql.Tx, t *entity.Transaction) error {
	query := `
		UPDATE transactions
		SET 
			account_id = $1, 
			category_id = $2, 
			type = $3, 
			amount = $4, 
			payee = $5,
			currency = $6, 
			transaction_date = $7, 
			description = $8, 
			receipt_url = $9, 
			location = $10, 
			updated_at = NOW()
		WHERE id = $11 AND user_id = $12
	`
	_, err := tx.ExecContext(ctx, query,
		t.AccountID, t.CategoryID, t.Type, t.Amount, t.Payee, t.Currency,
		t.TransactionDate, t.Description, t.ReceiptURL, t.Location,
		t.Id, t.UserID,
	)
	return err
}
