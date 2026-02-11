package storage

import (
	"context"
	"database/sql"

	"github.com/krakn/expense-management-backend-go/storage/datastore"
	"github.com/krakn/expense-management-backend-go/storage/entity"
	"github.com/krakn/expense-management-backend-go/storage/postgres"
	postgres_account "github.com/krakn/expense-management-backend-go/storage/postgres/account"
	postgres_category "github.com/krakn/expense-management-backend-go/storage/postgres/category"
	postgres_transaction "github.com/krakn/expense-management-backend-go/storage/postgres/transaction"
	postgres_user "github.com/krakn/expense-management-backend-go/storage/postgres/user"
)

type User interface {
	CreateUser(context.Context, datastore.Database, entity.User) (int64, error)
	GetUserByEmail(context.Context, datastore.Database, string) (entity.User, error)
	GetUserByID(context.Context, datastore.Database, int64) (entity.User, error)
	UpdateUser(context.Context, datastore.Database, entity.User) error
	DeleteUser(context.Context, datastore.Database, int64) error
}

type Category interface {
	CreateCategory(context.Context, datastore.Database, *entity.Category) (int64, error)
	GetCategoryByID(context.Context, datastore.Database, int64, int64) (*entity.Category, error)
	GetAllCategories(context.Context, datastore.Database, int64) (*[]entity.Category, error)
	UpdateCategory(context.Context, datastore.Database, *entity.Category, int64) error
	DeleteCategory(context.Context, datastore.Database, int64, int64) error
}

type Account interface {
	CreateAccount(context.Context, datastore.Database, *entity.Account) (int64, error)
	GetAccountByID(context.Context, datastore.Database, int64, int64) (*entity.Account, error)
	GetAllAccounts(context.Context, datastore.Database, int64) (*[]entity.Account, error)
	UpdateAccount(context.Context, datastore.Database, *entity.Account, int64) error
	DeleteAccount(context.Context, datastore.Database, int64, int64) error
}

type Transaction interface {
	CreateTransaction(context.Context, datastore.Database, *entity.Transaction) (int64, error)
	GetTransactionByID(context.Context, datastore.Database, int64, int64) (*entity.Transaction, error)
	GetAllTransactions(context.Context, datastore.Database, int64) ([]*entity.Transaction, error)
	UpdateTransaction(context.Context, datastore.Database, *entity.Transaction) error
	DeleteTransaction(context.Context, datastore.Database, int64, int64) error

	GetTransactionsByMonth(context.Context, datastore.Database, int64, int64, int64) ([]*entity.Transaction, error)
	GetTransactionsByCategory(context.Context, datastore.Database, int64, int64) ([]*entity.Transaction, error)
	GetTransactionsPaginated(context.Context, datastore.Database, int64, int64, int64) ([]*entity.Transaction, error)
	GetTransactionCount(context.Context, datastore.Database, int64) (int64, error)
}

type Storage struct {
	Connection interface {
		GetDb() *sql.DB
		Close() error
	}
	User        User
	Category    Category
	Account     Account
	Transaction Transaction
}

func NewStorage() (*Storage, error) {
	conn, err := postgres.CreateConfiguredPostgresConnection()
	if err != nil {
		return nil, err
	}
	return &Storage{
		Connection:  conn,
		User:        postgres_user.NewPostgresUserStorage(),
		Category:    postgres_category.NewPostgresCategoryStorage(),
		Account:     postgres_account.NewPostgresAccountStorage(),
		Transaction: postgres_transaction.NewPostgresTransactionStorage(),
	}, nil
}

func (s *Storage) WithTransaction(ctx context.Context, fn func(context.Context, datastore.Database) (any, error)) (any, error) {
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

func (s *Storage) Close() error {
	return s.Connection.Close()
}
