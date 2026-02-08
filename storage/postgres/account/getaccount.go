package postgres_account

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (s *PostgresAccountStorage) GetAccountByID(ctx context.Context, tx *sql.Tx, id int64, userId int64) (*entity.Account, error) {
	query := `
		SELECT 
			id, 
			name, 
			type, 
			currency, 
			current_balance, 
			bank_name, 
			account_number, 
			is_included_in_total, 
			user_id, 
			is_active, 
			created_at, 
			updated_at
		FROM 
			accounts
		WHERE 
			id = $1 AND user_id = $2
	`
	var acc entity.Account
	err := tx.QueryRowContext(ctx, query, id, userId).Scan(
		&acc.Id,
		&acc.Name,
		&acc.Type,
		&acc.Currency,
		&acc.CurrentBalance,
		&acc.BankName,
		&acc.AccountNumber,
		&acc.IsIncludedInTotal,
		&acc.UserID,
		&acc.IsActive,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (s *PostgresAccountStorage) GetAllAccounts(ctx context.Context, tx *sql.Tx, userID int64) (*[]entity.Account, error) {
	query := `
		SELECT 
			id, 
			name, 
			type, 
			currency, 
			current_balance, 
			bank_name, 
			account_number, 
			is_included_in_total, 
			user_id, 
			is_active, 
			created_at, 
			updated_at
		FROM 
			accounts
		WHERE 
			user_id = $1
	`
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []entity.Account
	for rows.Next() {
		var acc entity.Account
		err := rows.Scan(
			&acc.Id,
			&acc.Name,
			&acc.Type,
			&acc.Currency,
			&acc.CurrentBalance,
			&acc.BankName,
			&acc.AccountNumber,
			&acc.IsIncludedInTotal,
			&acc.UserID,
			&acc.IsActive,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return &accounts, nil
}
