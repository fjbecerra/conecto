package factories

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/sinks/codecs"
	"conecto/core/sinks/committers"
	"conecto/core/sinks/databases"
	"conecto/core/sinks/statestores"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Sink struct {
	Config SinkConfig
	FieldsSpecsConfigs map[string]FieldsSpecs
	Rand *rand.Rand
	StateStore statestores.StateStore
	DBConnection DBConnection
}

func NewSink(config SinkConfig, fieldsSpecsConfigs map[string]FieldsSpecs, rand *rand.Rand, stateStore statestores.StateStore, DBConnection DBConnection) *Sink{
	return &Sink{
		Config: config,
		FieldsSpecsConfigs: fieldsSpecsConfigs,
		Rand: rand,
		StateStore: stateStore,
		DBConnection: DBConnection,
	}
}

func (s * Sink) Build() engines.SinkEngine {
	var sink  sinks.Sink
	var commiter committers.Committer
	switch s.Config.Type {
		case PostgresSink:
			sink = buildPostgres(s.Config.SchemaConfig, s.FieldsSpecsConfigs)
			commiter = committers.NewDBCommiter(s.DBConnection.DB, sink, s.StateStore)			
			if s.Config.SchemaConfig.AutoCreate{
				s.createPostgresTable(s.Config.SchemaConfig.Table, s.FieldsSpecsConfigs[s.Config.SchemaConfig.FieldsSpecs], s.Config.SchemaConfig.Upsert)
			}
		case MemorySink:
			sink = buildSinkMemory()
			commiter = committers.NewMemoryCommiter(sink,s.StateStore)
		default:
			panic("unknown source type: " + s.Config.Type)
	}
	
	
	return engines.SinkEngine{
		Commiter: commiter,
		BatchSize: s.Config.BatchSize,
		MaxRetries: s.Config.Retry.MaxRetries,
		Backoff: time.Duration(s.Config.Retry.BackoffMS),
		MaxBackoff: time.Duration(s.Config.Retry.MaxBackoff),
		Rand: s.Rand,
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
		Upsert: buildUpsert(schemaConfig.Upsert),
		Adapter: &adapter,
	}
}

func buildSchema(tableName string, fieldsSpecs FieldsSpecs) databases.Schema {
	columns := []string{}
	for name, _ := range fieldsSpecs {		
		columns = append(columns, name)
	}
	metadata := []string{core.PipelineId}	
	return databases.Schema{
		Table: tableName,
		Columns: columns,
		Metadata: metadata,
	}
}

func buildUpsert(upsertConfig string)databases.Upsert {
	return databases.Upsert{
		ConflictColumns: strings.Split(upsertConfig, ","),
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

func (s * Sink) createPostgresTable(tableName string, specs FieldsSpecs, upsertConfig string){

	columns := []string{

		"    id SERIAL PRIMARY KEY",
		"    pipeline_id TEXT NOT NULL",
		"	 created_at TIMESTAMP DEFAULT NOW() NOT NULL",
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
		"    CONSTRAINT unique_constrains UNIQUE (" + upsertConfig + ")",
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
