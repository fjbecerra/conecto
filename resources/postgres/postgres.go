package postgres

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/statestores"
	"conecto/shared/clients"
	"conecto/shared/config"
)

type PostgresResource struct{
	client * clients.PostgresClient
	retryExecutor *retry.Executor
	stateStore statestores.StateStore
	config PostgresResourceConfig
}

func NewPostgresResource(
	client *clients.PostgresClient,
	retryExecutor *retry.Executor,
	stateStore statestores.StateStore,
	config PostgresResourceConfig,

) *PostgresResource{
	return &PostgresResource{
		client: client,
		retryExecutor: retryExecutor,
		stateStore: stateStore,
		config: config,
	}
}

func (p *PostgresResource) Close() error{
	return p.client.Close()
}

func (p *PostgresResource) Sink(cfg config.ConfigBytes, fieldsSpecs config.FieldsSpecs) engines.SinkCommiter {
    postgrestSinkConfig,_ := config.Unmarshal[PostgresSinkConfig](cfg, config.FormatJSON)
	postgresSink := PostgresSink{
		client: p.client,
		retryExecutor: p.retryExecutor,
		stateStore: p.stateStore,
		cfg: postgrestSinkConfig,
		fieldsSpecs: fieldsSpecs,
	}
	return CreatePostgresSink(postgresSink)
}

func (p *PostgresResource) Connector(cfg config.ConfigBytes) engines.ConnectorRunnable {
	return nil
}

func (p *PostgresResource) Transformers() []core.Transformer {
	return nil
}
