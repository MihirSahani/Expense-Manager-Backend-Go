package postgres

import "github.com/krakn/expense-management-backend-go/internal/utils"

type PostgresConfig struct {
	address  string
	MaxIdleConns int
	MaxOpenConns int
	MaxIdleTime  string
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		address: utils.GetEnv("POSTGRES_SERVER_ADDRESS", ""),
		MaxIdleConns: 10,
		MaxOpenConns: 25,
		MaxIdleTime:  "15m",
	}
}