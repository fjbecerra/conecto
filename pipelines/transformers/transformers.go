package transformers

import (
	"conecto/core"
	"conecto/shared/config"
)


func CreateTrasformer(fieldsSpecs config.FieldsSpecs) core.Transformer {
	tranformers := 	[]core.Transformer{}
	tranformers = append(tranformers, buildExtractor(fieldsSpecs))
	return &core.Chain{
		Steps : tranformers,
	}	
}

func buildExtractor(fieldsSpecs config.FieldsSpecs) core.Transformer {
	return &Extractor{
			Fields : Fields(fieldsSpecs),
			Selector: &GJSONSelector{},
	}    
}