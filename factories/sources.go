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
	NewMemory func() map[string]any
}


type Connection struct {
	OpenConnection OpenConnection
	Type SourcesType
}



type Connections map[string]Connection


func NewSource(SourcesConfig SourcesConfig) *Source {
	return &Source{
		SourcesConfig: SourcesConfig,
	}
}

func (d *Source) Build() Connections {
	connections:= make(Connections)
	for k, v := range d.SourcesConfig {
		openConnection := OpenConnection{}
		switch v.Type {
			case PostgresSource:
				openDb := func() *sql.DB{
					db, err := sql.Open("pgx", v.DSN)
					if err != nil {
						panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
					}
					return db
				}
				openConnection = OpenConnection{
					OpenDB: openDb,
				}	
		
			case MemorySource:
				memory := func() map[string]any{
					return make(map[string]any)
				}
				openConnection = OpenConnection{
					NewMemory: memory,
				}		
			default: 
				panic(fmt.Sprintf("Unknown source type: %s", v.Type))
		
		}
		connections[k]= Connection{
			OpenConnection: openConnection,
			Type: v.Type,
		}
		
	}
	return connections	
}