package factories

import (
	"database/sql"
	"fmt"
	"net/http"
)


type Source struct{
	SourcesConfig SourcesConfig
}


type OpenConnection struct {

	DB *sql.DB
	NewMemory  map[string]any
	httpClient *http.Client
}


type Connection struct {
	OpenConnection OpenConnection
	Type SourcesType
}

type Connections struct{
	connections map[string]Connection
} 

func (c *Connections) CloseAll() error {
	for _, connection := range c.connections {
		switch connection.Type {
		case PostgresSource:
			db := connection.OpenConnection.DB
			err := db.Close()
			if err != nil {
				return fmt.Errorf("failed to close Postgres connection: %w", err)
			}
		default:
			return fmt.Errorf("unknown source type: %s", connection.Type)
		}
	}
	return nil
}

func NewSource(SourcesConfig SourcesConfig) *Source {
	return &Source{
		SourcesConfig: SourcesConfig,
	}
}

func (d *Source) Build() Connections {
	connections:= make(map[string]Connection)
	for k, v := range d.SourcesConfig {
		openConnection := OpenConnection{}
		switch v.Type {
			case PostgresSource:
				db, err := sql.Open("pgx", v.DSN)
				if err != nil {
					panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
				}					
				openConnection = OpenConnection{
					DB: db,
				}	
		
			case MemorySource:
				memory := make(map[string]any)
				openConnection = OpenConnection{
					NewMemory: memory,
				}
			case HttpSource:
				httpClient := &http.Client{}
				openConnection = OpenConnection{
					httpClient: httpClient,
				}		
			default: 
				panic(fmt.Sprintf("Unknown source type: %s", v.Type))
		
		}
		connections[k]= Connection{
			OpenConnection: openConnection,
			Type: v.Type,
		}
		
	}
	return Connections{
		connections: connections,
	}	
}