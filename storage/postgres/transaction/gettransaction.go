package postgres_transaction

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresTransactionStorage) GetTransactionByID(ctx context.Context, tx *sql.Tx, id int64, userId int64) (*entity.Transaction, error) {
	query := `
		SELECT 
			id, user_id, account_id, category_id, type, amount, payee, currency, 
			transaction_date, description, receipt_url, location, created_at, updated_at
		FROM transactions
		WHERE id = $1 AND user_id = $2
	`
	var t entity.Transaction
	err := tx.QueryRowContext(ctx, query, id, userId).Scan(
		&t.Id, &t.UserID, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount, &t.Payee, &t.Currency,
		&t.TransactionDate, &t.Description, &t.ReceiptURL, &t.Location, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *PostgresTransactionStorage) GetAllTransactions(ctx context.Context, tx *sql.Tx, userId int64) ([]*entity.Transaction, error) {
	query := `
		SELECT 
			id, user_id, account_id, category_id, type, amount, payee, currency, 
			transaction_date, description, receipt_url, location, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY transaction_date DESC
	`
	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		err := rows.Scan(
			&t.Id, &t.UserID, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount, &t.Payee, &t.Currency,
			&t.TransactionDate, &t.Description, &t.ReceiptURL, &t.Location, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}
