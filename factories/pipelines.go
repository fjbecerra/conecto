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
	
	connector := NewConnector(config.ConnectorConfig, r).Build()
	transform := NewTransform(config.TransformersConfig, config.AdditionalConfig).Build()
	stateStore := NewStateStore(config.RuntimeConfig.StateStoreConfig).Build()
	sink :=NewSink(config.SinkConfig, config.AdditionalConfig, config.RuntimeConfig, r, stateStore).Build()
	
	return engines.Pipeline{
		ConnectorEngine: &connector,
		SinkEngine:   &sink,
        Transformer: transform,
        StateStore: stateStore,
	}
}