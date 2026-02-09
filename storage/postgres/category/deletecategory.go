package postgres_category

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/datastore"
)

func (p *PostgresCategoryStorage) DeleteCategory(ctx context.Context, tx datastore.Database, categoryID int64, userID int64) error {
	query := `
		DELETE FROM categories
		WHERE id = $1 AND user_id = $2
	`
	result, err := tx.ExecContext(ctx, query, categoryID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
