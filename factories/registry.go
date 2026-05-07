package factories

import "conecto/core/engines"

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
	return func() engines.PipelineRunner {
        config := LoadConfigPipeline(configPath)
		return BuildPipeline(config)
	}
}



