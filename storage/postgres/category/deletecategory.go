package postgres_category

import (
	"context"
	"database/sql"
)

func (p *PostgresCategoryStorage) DeleteCategory(ctx context.Context, tx *sql.Tx) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`
	result, err := tx.ExecContext(ctx, query)
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