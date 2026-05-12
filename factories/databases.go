package factories

import (
	"database/sql"
	"fmt"
)

type Database struct{
	DatabaseConfig DatabaseConfig
}

type DBConnection struct {
	DB      *sql.DB	
}

func NewDatabase(DatabaseConfig DatabaseConfig) *Database {
	return &Database{
		DatabaseConfig: DatabaseConfig,
	}
}

func (d *Database) Build() DBConnection {
	db, err := sql.Open("pgx", d.DatabaseConfig.DSN)
	if err != nil {
		panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
	}
	return DBConnection{
		DB: db,
	}	
}