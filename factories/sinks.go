package factories

import (
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/sinks/codecs"
	"conecto/core/sinks/rdbs"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Sink struct {
	Config SinkConfig
	AdditionalConfigs AdditionalConfig
	Rand *rand.Rand
}

func NewSink(config SinkConfig, additionalConfigs AdditionalConfig, rand *rand.Rand) *Sink{
	return &Sink{
		Config: config,
		AdditionalConfigs: additionalConfigs,
		Rand: rand,
	}
}

func (s * Sink) Build() engines.SinkEngine {
	var sink  sinks.Sink
	switch s.Config.Type {
		case Rdbs:
			sink = buildRdbs(s.Config.RDBSConfig, s.AdditionalConfigs)
		default:
			panic("unknown source type: " + s.Config.Type)
	}
	return engines.SinkEngine{
		Sink: sink,
		BatchSize: s.Config.BatchSize,
		MaxRetries: s.Config.Retry.MaxRetries,
		Backoff: time.Duration(s.Config.Retry.BackoffMS),
		MaxBackoff: time.Duration(s.Config.Retry.MaxBackoff),
		Rand: s.Rand,
	}
}

func buildRdbs(config RDBSConfig, additionalConfigs AdditionalConfig) *rdbs.Rdbs {
	switch config.DBType {
		case Postgres:
			return buildPostgres(config, additionalConfigs.FieldsConfig[config.Schema])
		default:
			panic("unkown rdbs type: " + config.DBType)
	}
	
}

func buildPostgres(config RDBSConfig, fieldsConfig FieldsConfig)*rdbs.Rdbs{
	jsonCodec := codecs.JSONCodec{}
	adapter := rdbs.PostgresAdapter{Codec: &jsonCodec}
	
	db, err := sql.Open("pgx", config.DSN)
	if err != nil {
		 panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
	}
	return &rdbs.Rdbs{
		DB: db,
		Schema: buildSchema(config.Table, fieldsConfig),
		Upsert: buildUpsert(config.Upsert),
		Adapter: &adapter,
	}
}

func buildSchema(tableName string, fieldsConfig FieldsConfig) rdbs.Schema {
	columns := []string{}
	for name, _ := range fieldsConfig {		
		columns = append(columns, name)
	}
	return rdbs.Schema{
		Table: tableName,
		Columns: columns,
	}
}

func buildUpsert(upsertConfig string)rdbs.Upsert {
	return rdbs.Upsert{
		ConflictColumns: strings.Split(upsertConfig, ","),
	}
}