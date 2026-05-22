package factories

import (
	"database/sql"
	"fmt"
)


type Source struct{
	SourcesConfig SourcesConfig
}


type OpenConnection struct {

	OpenDB func ()*sql.DB
}

type Connections map[string]OpenConnection


func NewSource(SourcesConfig SourcesConfig) *Source {
	return &Source{
		SourcesConfig: SourcesConfig,
	}
}

func (d *Source) Build() Connections {
	connections:= make(Connections)
	for k, v := range d.SourcesConfig {
		switch v.Type {
		case PostgresSource:
			openDb := func() *sql.DB{
				db, err := sql.Open("pgx", v.DSN)
				if err != nil {
					panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
				}
				return db
			}

			connections[k]= OpenConnection{
				OpenDB: openDb,
			}		
		}
	}
	return connections	
}