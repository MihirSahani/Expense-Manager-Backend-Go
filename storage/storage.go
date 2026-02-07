package storage

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/entity"
	"github.com/krakn/expense-management-backend-go/storage/postgres/user"
	"github.com/krakn/expense-management-backend-go/storage/postgres"
)

type Storage struct {
	Connection interface {
		GetDb() *sql.DB
		Close() error
	}

	User interface {
		CreateUser(context.Context, *sql.Tx, entity.User) (int64, error)
		GetUserByEmail(context.Context, *sql.Tx, string) (entity.User, error)
		GetUserByID(context.Context, *sql.Tx, int64) (entity.User, error)
		UpdateUser(context.Context, *sql.Tx, entity.User) error
		DeleteUser(context.Context, *sql.Tx, int64) error
	}

	Category interface {
		CreateCategory(context.Context, *sql.Tx, entity.Category) (int64, error)
		GetCategoryByID(context.Context, *sql.Tx) (entity.Category, error)
		GetAllCategories(context.Context, *sql.Tx) ([]entity.Category, error)
		UpdateCategory(context.Context, *sql.Tx, entity.Category) error
		DeleteCategory(context.Context, *sql.Tx) error
	}
}

func NewStorage() *Storage {
	conn, err := postgres.CreateConfiguredPostgresStorage()
	if err != nil {
		panic(err)
	}

	return &Storage{
		Connection: conn,
		User:       postgres_user.NewPostgresUserStorage(),
	}
}

func (s *Storage) WithTransaction(ctx context.Context, fn func(context.Context, *sql.Tx) (any, error)) (any,error) {
	tx, err := s.Connection.GetDb().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	data, err := fn(ctx, tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return data, nil
}