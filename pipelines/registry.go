package pipelines

type Registry[T any] struct {
    Factories map[string]T
}

func NewRegistryPipeline() *Registry[PipelineFactory] {
    return &Registry[PipelineFactory]{
        Factories: map[string]PipelineFactory{
            "fbAdInsight": func() PipelineRunner {
                return BuildPipeline("./testdata/fb_ad_insights/ad_insight_pipeline.json")
            },
        },
    }
}



