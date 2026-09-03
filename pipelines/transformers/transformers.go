package transformers

import (
	"conecto/core"
	"conecto/shared/config"
)


func CreateTrasformers(transformers []core.Transformer) core.Transformer {
	return &core.Chain{
		Steps : transformers,
	}	
}

func BuildExtractor(fieldsSpecs config.FieldsSpecs) core.Transformer {
	return &Extractor{
			Fields : Fields(fieldsSpecs),
			Selector: &GJSONSelector{},
	}    
}