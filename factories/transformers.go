package factories

import (
	"conecto/core"
	"conecto/transformers"
	
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

func (t * Transformer) Build() core.Transformer {
	tranformers := 	[]core.Transformer{}
	for _, tCfg := range t.Config {
		switch tCfg.Type {
			case TransformerExtractor : 
				fieldsConfig := t.FieldsSpecsConfig[tCfg.ExtractorConfig.Fields]
				tranformers = append(tranformers, buildExtractor(fieldsConfig))
			default:
					panic("unsupported source type")
		}
	}
	return &core.Chain{
		Steps : tranformers,
	}	
}

func buildExtractor(fieldsSpecs FieldsSpecs) core.Transformer {
	return &transformers.Extractor{
				Fields : transformers.Fields(fieldsSpecs),
				Selector: &transformers.GJSONSelector{},
	}    
}