package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	config *PostgresConfig
	Db     *sql.DB
}

func NewPostgresConnection(config *PostgresConfig) (*PostgresConnection, error) {
	
	db, err :=sql.Open("postgres", config.address)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	duration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	
	return &PostgresConnection{
		config: config,
		Db:     db,
	}, nil
}

func CreateConfiguredPostgresStorage() (*PostgresConnection, error) {
	return NewPostgresConnection(LoadPostgresConfig())
}

func (p *PostgresConnection) GetDb() *sql.DB {
	return p.Db
}

func (p *PostgresConnection) Close() error {
	p.Db.Close()
	return nil
}