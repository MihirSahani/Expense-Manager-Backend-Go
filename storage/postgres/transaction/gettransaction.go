package postgres_transaction

import (
	"context"

	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresTransactionStorage) GetTransactionByID(ctx context.Context, tx datastore.Database, id int64, userId int64) (*entity.Transaction, error) {
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

func (s *PostgresTransactionStorage) GetAllTransactions(ctx context.Context, tx datastore.Database, userId int64) ([]*entity.Transaction, error) {
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

func (s *PostgresTransactionStorage) GetTransactionsByMonth(ctx context.Context, tx datastore.Database, userId int64, month int64, year int64) ([]*entity.Transaction, error) {
	query := `
		SELECT id, account_id, category_id, type, amount, payee, currency, transaction_date, description, receipt_url, location, created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND EXTRACT(MONTH FROM transaction_date) = $2 AND EXTRACT(YEAR FROM transaction_date) = $3
		ORDER BY transaction_date DESC
	`

	var transactions []*entity.Transaction
	rows, err := tx.QueryContext(ctx, query, userId, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t entity.Transaction
		err := rows.Scan(
			&t.Id, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount, &t.Payee, &t.Currency,
			&t.TransactionDate, &t.Description, &t.ReceiptURL, &t.Location, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (s *PostgresTransactionStorage) GetTransactionsByCategory(ctx context.Context, tx datastore.Database, userId int64, categoryId int64) ([]*entity.Transaction, error) {
	query := `
		SELECT id, account_id, category_id, type, amount, payee, currency, transaction_date, description, receipt_url, location, created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND category_id = $2
		ORDER BY transaction_date DESC
	`
	var transactions []*entity.Transaction
	rows, err := tx.QueryContext(ctx, query, userId, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(
			&t.Id, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount, &t.Payee, &t.Currency, &t.TransactionDate, &t.Description, &t.ReceiptURL, &t.Location, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (s *PostgresTransactionStorage) GetTransactionsPaginated(ctx context.Context, tx datastore.Database, userId int64, page int64, pageSize int64) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, type, amount, payee, currency, transaction_date, description, receipt_url, location, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY transaction_date DESC
		LIMIT $2 OFFSET $3
	`
	var transactions []*entity.Transaction
	offset := (page - 1) * pageSize
	rows, err := tx.QueryContext(ctx, query, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(
			&t.Id, &t.UserID, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount, &t.Payee, &t.Currency, &t.TransactionDate, &t.Description, &t.ReceiptURL, &t.Location, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (s *PostgresTransactionStorage) GetTransactionCount(ctx context.Context, tx datastore.Database, userId int64) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM transactions
		WHERE user_id = $1
	`
	var count int64
	err := tx.QueryRowContext(ctx, query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
