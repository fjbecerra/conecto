package pipelines

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/pipelines/transformers"
	"conecto/resources"
	"conecto/shared/config"
	"fmt"
)

type Pipeline struct {
	Streams []engines.Stream
}

type PipelineRegistry struct {
	pipelineRegistry map[core.Provider]Pipeline
	resourceRegistry resources.ResourcesRegistry
}

func NewPipelineRegistry(resourceRegistry *resources.ResourcesRegistry)PipelineRegistry{
	return PipelineRegistry{
		pipelineRegistry: make(map[core.Provider]Pipeline),
		resourceRegistry: *resourceRegistry,
	}
}

func (p *PipelineRegistry) Register(pipelines []config.Pipeline) error {

	for _, pipe := range pipelines {
		streamsCfg, error := config.LoadPipelineConfig[config.Streams](pipe.Path)
		if(error!=nil){
			fmt.Errorf("pipeline path not found")
		}
		streams := p.createStreams(streamsCfg)
		
		p.pipelineRegistry[core.Provider(pipe.Provider)]=Pipeline{Streams: streams}
	}

	return nil	
	
}

func (p *PipelineRegistry)  createStreams(cfg config.Streams) []engines.Stream{
	streams := []engines.Stream{}
	for _, st := range cfg.Streams{
		
		inbound := p.resourceRegistry.Get(resources.ResourceName(st.Inbound.Resource))
		connector := inbound.Connector(st.Inbound.Config)
		outbound:= p.resourceRegistry.Get(resources.ResourceName(st.Outbound.Resource))
		sink:= outbound.Sink(st.Outbound.Config, st.FieldsSpecs)
		tfs := 	[]core.Transformer{}
		tfs = append(tfs, transformers.BuildExtractor(st.FieldsSpecs))
		inboundTransformers := inbound.Transformers()
		tfs = append(tfs, inboundTransformers...)
		transformer:= transformers.CreateTrasformers(tfs)		
		engine:= engines.Engine {
			ConnectorRunnable: connector,
			Transformer: transformer,
			SinkCommiter: sink,
		}

		stream := engines.Stream {
			Name: st.Name,
			Engine: &engine,
			StateStore: p.resourceRegistry.StateStore,
		}
		
		streams = append(streams, stream)
	}

	return streams

}

func (p *PipelineRegistry) Get(provider core.Provider) Pipeline{
	return p.pipelineRegistry[provider]
}