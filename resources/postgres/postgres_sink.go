package postgres

import (
	"conecto/core/codecs"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/statestores"
	"conecto/resources/base/db"
	"conecto/shared/clients"
	"conecto/shared/config"
	"database/sql"
	"fmt"
	"strings"
)

type PostgresSink struct{
	client * clients.PostgresClient
	retryExecutor *retry.Executor
	stateStore statestores.StateStore	
	cfg PostgresSinkConfig
	fieldsSpecs config.FieldsSpecs
}

func CreatePostgresSink(postgresSink PostgresSink) engines.SinkCommiter{
	tableName := postgresSink.cfg.Table
	autoCreate:= postgresSink.cfg.AutoCreate

	sinker := createPostgres(postgresSink.client, tableName, postgresSink.fieldsSpecs, autoCreate,)
	return &engines.SinkEngine{
		Sinker: *sinker,
		SinkRetry: *postgresSink.retryExecutor,
		StateStore: postgresSink.stateStore,
	}
}


func createPostgres(
	client *clients.PostgresClient,
	tableName string,
	fieldSpecs config.FieldsSpecs,
	autoCreate bool,

) *engines.Sinker {
	database := client.Get()
	sink := buildPostgres(database, tableName, fieldSpecs)	
	if (autoCreate){
		createPostgresTable(database, tableName, fieldSpecs)
	}
		
	sqlExecutor := db.SQLExecutor{
		OpenTransaction: func() *sql.Tx {
			tx, err := database.Begin()
			if err != nil {
				panic(err)
			}
			return tx
		},
	}		
	return &engines.Sinker{
		Sink: sink,				
		Executor: &sqlExecutor,
	}			 
}

func buildPostgres(database *sql.DB, tableName string, fieldsSpecs config.FieldsSpecs)*db.Rdbs{
	jsonCodec := codecs.JSONCodec{}
	adapter := PostgresAdapter{Codec: &jsonCodec, Upsert: true }
	
	return &db.Rdbs{
		Schema: buildSchema(database, tableName, fieldsSpecs),
		Adapter: &adapter,
	}
}

func buildSchema(database *sql.DB, tableName string, fieldsSpecs config.FieldsSpecs) db.Schema {
	columns := []string{}
	for name, _ := range fieldsSpecs {		
		columns = append(columns, name)
	}
	metadata := []db.Metadata{
		{
			Column: "__stream_id",
			Unique: true,
		},
		{
			Column: "__event_id",
			Unique: true,
		},
		{
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

func createPostgresTable(database *sql.DB, tableName string, specs config.FieldsSpecs){

	columns := []string{

		"    __id SERIAL PRIMARY KEY",
		"	 __event_id TEXT NOT NULL",
		"    __stream_id TEXT NOT NULL",
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

	_, err := database.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("table created or already exists")
}