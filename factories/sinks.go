package factories

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sinks/rdbs"
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Sink struct {
	Config core.SinkConfig
	AdditionalConfigs core.AdditionalConfig
}

func NewSink(config core.SinkConfig, additionalConfigs core.AdditionalConfig) *Sink{
	return &Sink{
		Config: config,
		AdditionalConfigs: additionalConfigs,
	}
}

func (sink * Sink) Build() sinks.Sink[core.Record] {
	switch sink.Config.Type {
		case core.Rdbs:
			return buildRdbs(sink.Config.RDBSConfig, sink.AdditionalConfigs)
		default:
			panic("unknown source type: " + sink.Config.Type)
	}
}

func buildRdbs(config core.RDBSConfig, additionalConfigs core.AdditionalConfig) *rdbs.Rdbs {
	switch config.DBType {
		case core.Postgres:
			return buildPostgres(config, additionalConfigs[config.Schema])
		default:
			panic("unkown rdbs type: " + config.DBType)
	}
	
}

func buildPostgres(config core.RDBSConfig, fieldsConfig core.FieldsConfig)*rdbs.Rdbs{
	adapter := rdbs.PostgresAdapter{}
	
	
	db, err := sql.Open("pgx", config.DSN)
	if err != nil {
		 panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
	}
	return &rdbs.Rdbs{
		DB: db,
		Schema: buildSchema(config.Table, fieldsConfig),
		Upsert: buildUpsert(config.Upsert),
		Adapter: &adapter,
		BatchSize: config.BatchSize,
	}
}

func buildSchema(tableName string, fieldsConfig core.FieldsConfig) rdbs.Schema {
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