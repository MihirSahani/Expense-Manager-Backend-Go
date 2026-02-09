package postgres_transaction

type PostgresTransactionStorage struct{}

func NewPostgresTransactionStorage() *PostgresTransactionStorage {
	return &PostgresTransactionStorage{}
}
