package postgres

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
)

type PostgresUserStorage struct{}

func NewPostgresUserStorage() *PostgresUserStorage {
	return &PostgresUserStorage{}
}

func (p *PostgresUserStorage) CreateUser(ctx context.Context, tx *sql.Tx, user entity.User) (int64, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostgresUserStorage) GetUserByEmail(ctx context.Context, tx *sql.Tx, email string) (entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password
		FROM users
		WHERE email = $1
	`

	var user entity.User
	err := tx.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (p *PostgresUserStorage) GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (entity.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password
		FROM users
		WHERE id = $1
	`
	var user entity.User
	err := tx.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (p *PostgresUserStorage) UpdateUser(ctx context.Context, tx *sql.Tx, user entity.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password = $4
		WHERE id = $5
	`
	_, err := tx.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserStorage) DeleteUser(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := tx.ExecContext(ctx, query, id)
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
