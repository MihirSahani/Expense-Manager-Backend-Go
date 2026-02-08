package postgres_account

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresAccountStorage) CreateAccount(ctx context.Context, tx *sql.Tx, acc *entity.Account) (int64, error) {
	query := `
		INSERT INTO 
			accounts (name, type, currency, current_balance, bank_name, account_number, is_included_in_total, user_id, is_active)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query,
		acc.Name,
		acc.Type,
		acc.Currency,
		acc.CurrentBalance,
		acc.BankName,
		acc.AccountNumber,
		acc.IsIncludedInTotal,
		acc.UserID,
		acc.IsActive,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}
