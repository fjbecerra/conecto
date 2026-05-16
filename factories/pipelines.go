package factories

import (
	"conecto/core/engines"
	"math/rand/v2"
	"time"
)



func BuildPipeline(config ConfigPipeline) engines.Pipeline{
	seed := time.Now().UnixNano()

	r := rand.New(rand.NewPCG(
		uint64(seed),
		uint64(seed>>1),
	))

	connection:= NewDatabase(config.DatabaseConfig).Build()
	connector := NewConnector(config.ConnectorConfig, r.Float64, connection).Build()
	transform := NewTransform(config.TransformersConfig, config.FieldsSpecsConfig, config.RuntimeConfig).Build()
	stateStore := NewStateStore(config.RuntimeConfig.StateStoreConfig, connection).Build()
	sink := NewSink(config.SinkConfig, config.FieldsSpecsConfig, r.Float64, stateStore, connection).Build()
	
	return engines.Pipeline{
		ConnectorEngine: &connector,
		CommitStrategy:  sink,
        Transformer: transform,
        StateStore: stateStore,
	}
}