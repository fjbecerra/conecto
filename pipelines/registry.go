package pipelines

import (
    "time"
    "conecto/core"
)

type Registry[T any] struct {
    Factories map[string]T
}

func NewRegistryPipeline() *Registry[PipelineFactory] {
	return &Registry[PipelineFactory]{
		Factories: map[string]PipelineFactory{
            "fbAdInsight": BuildFactory("./testdata/fb_ad_insights/ad_insight_pipeline.json"),
        },
	}
}

func BuildFactory(configPath string) PipelineFactory {
	return func() PipelineRunner {
        config := core.LoadConfigPipeline(configPath)
		p := BuildPipeline(config)

		if !config.AdditionalConfig.RetryConfig.Enabled {
			return p
		}

		return &RetryRunner{
			inner:      p,
			maxRetries: config.AdditionalConfig.RetryConfig.MaxRetries,
			backoff:    time.Duration(config.AdditionalConfig.RetryConfig.BackoffMS) * time.Millisecond,
		}
	}
}




