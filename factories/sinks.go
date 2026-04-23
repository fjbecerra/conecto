package factories

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sinks/rdbs"
	"database/sql"
	"fmt"
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
		case core.SinkMemory:
			return sink.buildMemorySink()
		case core.Rdbs:
			return buildRdbs(sink.Config.RDBSConfig, sink.AdditionalConfigs)
		default:
			panic("unknown source type: " + sink.Config.Type)
	}
}

func (sink *Sink) buildMemorySink() *sinks.SinkMemory[core.Record] {
	return sinks.NewMemorySink[core.Record]()
}

func buildRdbs(config core.RDBSConfig, additionalConfigs core.AdditionalConfig) *rdbs.Rdbs {
	switch config.DBType {
		case core.Postgres:
			return buildPostgres(config, additionalConfigs[config.Schema].FieldsConfig)
		default:
			panic("unkown rdbs type: " + config.DBType)
	}
	
}

func buildPostgres(config core.RDBSConfig, fieldsConfig core.FieldsConfig)*rdbs.Rdbs{
	adapter := rdbs.PostgresAdapter{}
	
	db, err := sql.Open(adapter.DriverName(), config.DSN)
	if err != nil {
		 panic(fmt.Sprintf("cannot open connection, %s", err.Error()))
	}
	return &rdbs.Rdbs{
		DB: db,
		Table: config.Table,
		Schema: buildSchema(fieldsConfig),
		Adapter: &adapter,
		BatchSize: config.BatchSize,
	}

}

func buildSchema(fieldsConfig core.FieldsConfig) rdbs.Schema {
	fields := []rdbs.Field{}
	for name,cfg := range fieldsConfig {
		field := rdbs.Field {
			Name: name,
			Default: cfg.Default,
		}
		fields = append(fields, field)
	}
	return rdbs.Schema{Fields: fields}
}