package clients

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	Dsn string `json:"dsn"`
}

type PostgresClient struct{
	sql *sql.DB
}

func CreatePostgresClient(config PostgresConfig) *PostgresClient{
	return newPostgresClient(config)
}

func newPostgresClient(config PostgresConfig) *PostgresClient{
	db, err := sql.Open("pgx", config.Dsn)
	if err!=nil {
		panic(fmt.Sprintf("%s for dsn %s", err, config.Dsn))
	}
	return &PostgresClient{
		sql: db,
	}		
}

func(p * PostgresClient) Get() *sql.DB {
	return 	p.sql		
}

func(p * PostgresClient) Close() error {
	return p.sql.Close()
}