package factories

import (
	"conecto/core"
	"conecto/core/transformers"
)

type Transformer struct {
	Config []core.TransformerConfig
	AdditinalConfigs core.AdditionalConfig

}

func NewTransform(config []core.TransformerConfig, additinalConfigs core.AdditionalConfig) *Transformer {
	return &Transformer {
		Config: config,
		AdditinalConfigs: additinalConfigs,
	}
}

func (t * Transformer) Build() transformers.Transformer {
	tranformers := 	[]transformers.Transformer{}
	for _, tCfg := range t.Config {
		switch tCfg.Type {
			case core.TransformerExtractor : 
				fieldsConfig := t.AdditinalConfigs.FieldsConfig[tCfg.ExtractorConfig.Fields]
				tranformers = append(tranformers, buildExtractor(fieldsConfig))
			default:
					panic("unsupported source type")
			}
		}
	return &transformers.Chain{
		Steps : tranformers,
	}	
}

func buildExtractor(fieldsConfig core.FieldsConfig) transformers.Transformer {
	return &transformers.Extractor{
				Fields : transformers.Fields(fieldsConfig),
				Selector: &transformers.GJSONSelector{},
	}    
}