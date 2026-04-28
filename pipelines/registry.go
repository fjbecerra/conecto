package pipelines

import (
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
		return BuildPipeline(config)
	}
}




