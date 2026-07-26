package factories

import "conecto/pipelines"


type Conecto struct{
	
}

func Build(conectoConfig ConectoConfig) Conecto{
	random:= &RandomImpl{}
	connections:= NewSource(conectoConfig.SourcesConfig).Build()
	stateStore := NewStateStore(conectoConfig.RuntimeConfig.StateStoreConfig, connections).Build()
	pipelineRegiter := pipelines.NewRegistry()
	for _, path := range conectoConfig.PipelineRegisterConfig {
		pipelineConfig, error := LoadConfig[PipelineConfig](path)
		if(error!=nil){
			panic("path not found")
		}
		
		
		pipeline:= NewPipeline(connections, random, stateStore, ,pipelineConfig)
		pipelineRegiter.Register(pipeline)
	}






	return Conecto{
		
	}
}

func builStateStore() {

}

