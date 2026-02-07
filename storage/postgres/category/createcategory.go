package postgres_category

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (p *PostgresCategoryStorage) CreateCategory(ctx context.Context, tx *sql.Tx, category entity.Category) (int64, error) {
	query := `
		INSERT INTO categories (name, type, color, desc, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query, category.Name, category.Type, category.Color, category.Desc, category.UserID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}