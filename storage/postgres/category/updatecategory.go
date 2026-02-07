package postgres_category

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (p *PostgresCategoryStorage) UpdateCategory(ctx context.Context, tx *sql.Tx, category entity.Category) error {
	query := `
		UPDATE categories
		SET name = $1, type = $2, color = $3, desc = $4, updated_at = NOW() 
		WHERE id = $5
	`
	result, err := tx.ExecContext(ctx, query, category.Name, category.Type, category.Color, category.Desc, category.Id)
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