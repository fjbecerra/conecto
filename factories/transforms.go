package factories

import (
	"conecto/core"
	"conecto/core/sources"
	"conecto/core/transforms"
	"encoding/json"
)

type Transform struct {
	Source any
	Config []core.TransformConfig
	AdditinalConfigs core.AdditionalConfig

}

func NewTransform(source any, config []core.TransformConfig, additinalConfigs core.AdditionalConfig) *Transform {
	return &Transform {
		Source: source,
		Config: config,
		AdditinalConfigs: additinalConfigs,
	}
}

func (transform *Transform) Build() sources.Source[core.Record] {
	var current sources.Source[core.Record]

	for _, t := range transform.Config {
		switch src := transform.Source.(type) {
			case sources.Source[json.RawMessage]:
				var fn  core.Transform[json.RawMessage,core.Record]
				switch t.Type {
					case core.TransformExtractor :						
						fieldsConfig := transform.AdditinalConfigs.FieldsConfig[t.ExtractorConfig.Fields]
						fn = buildExtractor(fieldsConfig)
					default:
						panic("unsupported transform")
				}
				current = &core.MapSource[json.RawMessage, core.Record]{
					Upstream: src,
					MapFn: fn,
					Workers: t.Workers,
				}
			default:
				panic("unsupported source type")
		}
	}
	return current
}

func buildExtractor(fieldsConfig core.FieldsConfig) core.Transform[json.RawMessage, core.Record] {
	extractor :=transforms.Extractor{
					Fields : transforms.Fields(fieldsConfig),
				}
    return  func (in json.RawMessage) core.Record {
		record, _ := extractor.Extract(in)
		return record
	}  
}