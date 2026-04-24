package pipelines

type Registry[T any] struct {
    Factories map[string]T
}

func NewRegistryPipeline() *Registry[PipelineFactory] {
    return &Registry[PipelineFactory]{
        Factories: map[string]PipelineFactory{
            "fbAdInsight": func() PipelineRunner {
                return BuildPipeline("../config/facebook_ad_insight_pipeline.json")
            },
             "pokeApi": func() PipelineRunner {
                return BuildPipeline("./testdata/poke_api/poke_api_pipeline.json")
            },
        },
    }
}



