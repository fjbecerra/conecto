package factories

import (
	"conecto/core/engines"
	"math/rand/v2"
	"time"
)

func BuildPipeline(config ConfigPipeline) engines.PipelineRunner{
	seed := time.Now().UnixNano()

	r := rand.New(rand.NewPCG(
		uint64(seed),
		uint64(seed>>1),
	))
	
	connector := NewConnector(config.ConnectorConfig, r).Build()
	transform := NewTransform(config.TransformersConfig, config.AdditionalConfig).Build()
	sink :=NewSink(config.SinkConfig, config.AdditionalConfig, r).Build()
	settings := engines.Settings {
		BufferSize: config.AdditionalConfig.BufferSize,
	}

	return &engines.Pipeline{
		ConnectorEngine: &connector,
		SinkEngine:   &sink,
        Transformer: transform,
        Settings: settings,
	}
}


type PipelineFactory func() engines.PipelineRunner