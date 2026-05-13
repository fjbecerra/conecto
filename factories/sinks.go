package factories

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/sinks"
	"conecto/core/sinks/codecs"
	"conecto/core/sinks/databases"
	"conecto/core/statestores"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Sink struct {
	Config SinkConfig
	FieldsSpecsConfigs map[string]FieldsSpecs
	RandFn func() float64
	StateStore statestores.StateStore
	DBConnection DBConnection
}

func NewSink(config SinkConfig, fieldsSpecsConfigs map[string]FieldsSpecs, randFn func() float64, stateStore statestores.StateStore, DBConnection DBConnection) *Sink{
	return &Sink{
		Config: config,
		FieldsSpecsConfigs: fieldsSpecsConfigs,
		RandFn: randFn,
		StateStore: stateStore,
		DBConnection: DBConnection,
	}
}

func (s * Sink) Build() engines.CommitStrategy {
	var sink  sinks.Sink
	retryPolicy:= retry.Policy{
		MaxRetries: s.Config.Retry.MaxRetries,
		InitialBackoff: time.Duration(s.Config.Retry.BackoffMS),
		MaxBackoff: time.Duration(s.Config.Retry.MaxBackoff),
		Jitter: true,
	}
	retryExecutor := retry.Executor {
		Policy: retryPolicy,
		Rand: s.RandFn,
	}
	switch s.Config.Type {
		case PostgresSink:
			sink = buildPostgres(s.Config.SchemaConfig, s.FieldsSpecsConfigs)			
			if s.Config.SchemaConfig.AutoCreate{
				s.createPostgresTable(s.Config.SchemaConfig.Table, s.FieldsSpecsConfigs[s.Config.SchemaConfig.FieldsSpecs])
			}			
			return &engines.Transactional{
				DB: s.DBConnection.DB,
				Sink: sink,
				StateStore: s.StateStore,
				Retry: retryExecutor,
			}	
		case MemorySink:
			
			sink = buildSinkMemory()
			return &engines.AtLeastOnceCommitStrategy{
				SinkRetry: retryExecutor,
				Sink: sink,
				StateStore: s.StateStore,
				StateStoreRetry: retryExecutor,
				Executor: &commands.MemoryCommandExecutor{
					Store: &[]core.Event{},
				},
			}
		default:
			panic("unknown source type: " + s.Config.Type)
	}
	
}

func  buildSinkMemory() *sinks.SinkMemory {	 
	mstore:= []map[string]interface{}{}
	 jsonCodec := codecs.JSONCodec{}	
	return sinks.NewMemorySink(mstore, &jsonCodec)
}


func buildPostgres(schemaConfig SchemaConfig, fieldsSpecsConfig map[string]FieldsSpecs)*databases.Rdbs{
	jsonCodec := codecs.JSONCodec{}
	adapter := databases.PostgresAdapter{Codec: &jsonCodec}
	
	return &databases.Rdbs{
		Schema: buildSchema(schemaConfig.Table, fieldsSpecsConfig[schemaConfig.FieldsSpecs]),
		Adapter: &adapter,
	}
}

func buildSchema(tableName string, fieldsSpecs FieldsSpecs) databases.Schema {
	columns := []string{}
	for name, _ := range fieldsSpecs {		
		columns = append(columns, name)
	}
	metadata := []databases.Metadata{
		databases.Metadata{
			Column: "__pipeline_id",
			Unique: true,
		},
		databases.Metadata{
			Column: "__event_id",
			Unique: true,
		},
		databases.Metadata{
			Column: "__ingested_at",
			Unique: false,
		},
	}	
	return databases.Schema{
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

func (s * Sink) createPostgresTable(tableName string, specs FieldsSpecs){

	columns := []string{

		"    id SERIAL PRIMARY KEY",
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
		"    CONSTRAINT unique_constrains UNIQUE (__event_id, __pipeline_id)",
	)

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	%s
	);
	`, tableName, strings.Join(columns, ",\n"))

	_, err := s.DBConnection.DB.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}
