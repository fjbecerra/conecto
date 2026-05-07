package factories

import (
	"conecto/core/transformers"
	
)

type Transformer struct {
	Config []TransformerConfig
	AdditinalConfigs AdditionalConfig

}

func NewTransform(config []TransformerConfig, additinalConfigs AdditionalConfig) *Transformer {
	return &Transformer {
		Config: config,
		AdditinalConfigs: additinalConfigs,
	}
}

func (t * Transformer) Build() transformers.Transformer {
	tranformers := 	[]transformers.Transformer{}
	for _, tCfg := range t.Config {
		switch tCfg.Type {
			case TransformerExtractor : 
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

func buildExtractor(fieldsConfig FieldsConfig) transformers.Transformer {
	return &transformers.Extractor{
				Fields : transformers.Fields(fieldsConfig),
				Selector: &transformers.GJSONSelector{},
	}    
}