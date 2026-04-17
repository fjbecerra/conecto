package pipelines

import "conecto/core"

type Registry struct {
    Factories map[string]PipelineFactory
}

func NewRegistry() *Registry {
    return &Registry{
        Factories: map[string]PipelineFactory{
            "facebookAdInsight": func() Pipeline[core.Record] {
                return *NewPipeline("./schemas/facebook_ad_insight.json")
            },
        },
    }
}