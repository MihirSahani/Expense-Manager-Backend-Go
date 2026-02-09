package postgres_user

type PostgresUserStorage struct{}

func NewPostgresUserStorage() *PostgresUserStorage {
	return &PostgresUserStorage{}
}
