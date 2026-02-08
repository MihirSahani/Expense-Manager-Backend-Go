package postgres_account

type PostgresAccountStorage struct{}

func NewPostgresAccountStorage() *PostgresAccountStorage {
	return &PostgresAccountStorage{}
}