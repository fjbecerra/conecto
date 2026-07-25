package factories

import (
	"conecto/core"
	"conecto/core/codecs"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/statestores"
	"conecto/sinks"
	"conecto/sinks/base/db"
	"conecto/sinks/memory"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Sink struct {
	Config SinkConfig
	FieldsSpecsConfigs map[string]FieldsSpecs
	random retry.Random
	StateStore statestores.StateStore
	connections Connections
	destinationConfig DestinationConfig
}

func NewSink(config SinkConfig, fieldsSpecsConfigs map[string]FieldsSpecs, random retry.Random, stateStore statestores.StateStore, connections Connections, destinationConfig DestinationConfig) *Sink{
	return &Sink{
		Config: config,
		FieldsSpecsConfigs: fieldsSpecsConfigs,
		random: random,
		StateStore: stateStore,
		connections: connections,
		destinationConfig: destinationConfig,
	}
}

func (s * Sink) Build() engines.SinkCommiter {
	var sink  core.Sink
	retryPolicy:= retry.Policy{
		MaxRetries: s.Config.Retry.MaxRetries,
		InitialBackoff: time.Duration(s.Config.Retry.BackoffMS),
		MaxBackoff: time.Duration(s.Config.Retry.MaxBackoff),
		Jitter: true,
	}
	retryExecutor := retry.Executor {
		Policy: retryPolicy,
		Random: s.random,
	}
	switch s.Config.Type {
		case PostgresSink:
			sink = buildPostgres(s.destinationConfig, s.FieldsSpecsConfigs)	
			connection:= s.connections[s.Config.Source].OpenDB()	
			if s.Config.SchemaConfig.AutoCreate{
				createPostgresTable(s.destinationConfig.Name, s.FieldsSpecsConfigs[*s.destinationConfig.Schema],connection)
			}
				
			sqlExecutor := db.SQLExecutor{
				OpenTransaction: func() *sql.Tx {
					tx, err := connection.Begin()
					if err != nil {
						panic(err)
					}
					return tx
				},
			}		
			return &engines.SinkEngine{
				Sink: sink,				
				SinkRetry: retryExecutor,
				StateStore: s.StateStore,
				StateStoreRetry: retryExecutor,
				Executor: &sqlExecutor,
			}			 
		case MemorySink:
			
			sink = buildSinkMemory()
			
			return &engines.SinkEngine{
				SinkRetry: retryExecutor,
				Sink: sink,
				StateStore: s.StateStore,
				StateStoreRetry: retryExecutor,
				Executor: &memory.MemoryExecutor{
					Store: &[]core.Event{},
				},
			}
		default:
			panic("unknown source type: " + s.Config.Type)
	}
	
}

func  buildSinkMemory() *memory.SinkMemory {	 
	mstore:= []map[string]interface{}{}
	 jsonCodec := codecs.JSONCodec{}	
	return memory.NewMemorySink(mstore, &jsonCodec)
}


func buildPostgres(destinationConfig DestinationConfig, fieldsSpecsConfig map[string]FieldsSpecs)*db.Rdbs{
	jsonCodec := codecs.JSONCodec{}
	adapter := sinks.Postgres{Codec: &jsonCodec, Upsert: true }
	
	return &db.Rdbs{
		Schema: buildSchema(destinationConfig.Name, fieldsSpecsConfig[*destinationConfig.Schema]),
		Adapter: &adapter,
	}
}

func buildSchema(tableName string, fieldsSpecs FieldsSpecs) db.Schema {
	columns := []string{}
	for name, _ := range fieldsSpecs {		
		columns = append(columns, name)
	}
	metadata := []db.Metadata{
		db.Metadata{
			Column: "__pipeline_id",
			Unique: true,
		},
		db.Metadata{
			Column: "__event_id",
			Unique: true,
		},
		db.Metadata{
			Column: "__ingested_at",
			Unique: false,
		},
	}	
	return db.Schema{
		Table: tableName,
		Columns: columns,
		Metadata: metadata,
	}
}

var pgTypes = map[string]string{
	"float":  "DOUBLE PRECISION",
	"int":    "INTEGER",
	"time":   "TIMESTAMP",
	"string": "TEXT",
	"bool":   "BOOLEAN",
	"int64":   "BIGINT",
}

func createPostgresTable(tableName string, specs FieldsSpecs, db *sql.DB){

	columns := []string{

		"    __id SERIAL PRIMARY KEY",
		"	 __event_id TEXT NOT NULL",
		"    __pipeline_id TEXT NOT NULL",
		"	 __ingested_at TIMESTAMP NOT NULL",
		"	 __created_at TIMESTAMP DEFAULT NOW() NOT NULL",
	}

	for name, field := range specs {

		pgType, ok := pgTypes[field.Type]
		if !ok {
			panic("unsupported type: " + field.Type)
		}

		col := fmt.Sprintf("%s %s", name, pgType)

		if field.Required {
			col += " NOT NULL"
		}

		if field.Default != nil {

			switch v := field.Default.(type) {

			case string:

				// postgres function like NOW()
				if strings.HasSuffix(v, "()") {
					col += fmt.Sprintf(" DEFAULT %s", v)
				} else {
					col += fmt.Sprintf(" DEFAULT '%s'", v)
				}

			case bool:
				col += fmt.Sprintf(" DEFAULT %t", v)

			default:
				col += fmt.Sprintf(" DEFAULT %v", v)
			}
		}

		columns = append(columns, "    "+col)
	}

	// add unique constraint
	columns = append(
		columns,
		"CONSTRAINT unique_constrains_" + tableName +  " UNIQUE (__event_id)",
	)

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	%s
	);
	`, tableName, strings.Join(columns, ",\n"))

	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
