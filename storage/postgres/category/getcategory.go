package postgres_category

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

func (p *PostgresCategoryStorage) GetCategoryByID(ctx context.Context, tx *sql.Tx) (*entity.Category, error) {
	query := `
		SELECT id, name, type, color, desc, user_id, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	var category entity.Category
	err := tx.QueryRowContext(ctx, query, category.Id).Scan(&category.Id, &category.Name, &category.Type, &category.Color, &category.Desc, &category.UserID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (p *PostgresCategoryStorage) GetAllCategories(ctx context.Context, tx *sql.Tx, userID int64) (*[]entity.Category, error) {
	query := `
		SELECT id, name, type, color, desc, user_id, created_at, updated_at
		FROM categories
		WHERE user_id = $1
	`
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Type, &category.Color, &category.Desc, &category.UserID, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return &categories, nil
}