package postgres_account

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresAccountStorage) UpdateAccount(ctx context.Context, tx *sql.Tx, acc *entity.Account, userId int64) error {
	query := `
		UPDATE 
			accounts
		SET 
			name = $1, 
			type = $2, 
			currency = $3, 
			current_balance = $4, 
			bank_name = $5, 
			account_number = $6, 
			is_included_in_total = $7, 
			is_active = $8, 
			updated_at = NOW()
		WHERE 
			id = $9 AND user_id = $10;
	`
	_, err := tx.ExecContext(ctx, query, acc.Name, acc.Type, acc.Currency, acc.CurrentBalance, acc.BankName, acc.AccountNumber, acc.IsIncludedInTotal, acc.IsActive, acc.Id, userId)
	return err
}
