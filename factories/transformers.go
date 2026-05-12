package factories

import (
	"conecto/core/transformers"
	
)

type Transformer struct {
	Config []TransformerConfig
	FieldsSpecsConfig map[string]FieldsSpecs
	RuntimeConfig RuntimeConfig

}

func NewTransform(config []TransformerConfig, fieldsSpecsConfig map[string]FieldsSpecs, runtimeConfig RuntimeConfig) *Transformer {
	return &Transformer {
		Config: config,
		FieldsSpecsConfig: fieldsSpecsConfig,
		RuntimeConfig: runtimeConfig,
 	}
}

func (t * Transformer) Build() transformers.Transformer {
	tranformers := 	[]transformers.Transformer{}
	for _, tCfg := range t.Config {
		switch tCfg.Type {
			case TransformerExtractor : 
				fieldsConfig := t.FieldsSpecsConfig[tCfg.ExtractorConfig.Fields]
				tranformers = append(tranformers, buildExtractor(fieldsConfig))
			default:
					panic("unsupported source type")
		}
	}
	//append internal transformers
	tranformers = append(tranformers, buildTenantIdEnricher(t.RuntimeConfig.PipelineId))
	return &transformers.Chain{
		Steps : tranformers,
	}	
}

func buildExtractor(fieldsSpecs FieldsSpecs) transformers.Transformer {
	return &transformers.Extractor{
				Fields : transformers.Fields(fieldsSpecs),
				Selector: &transformers.GJSONSelector{},
	}    
}

func buildTenantIdEnricher(pipelineId string) transformers.Transformer {
	return &transformers.PipelineIdEnricher{
				PipelineId: pipelineId,
	}    
}