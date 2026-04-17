package extractors

type FieldsConfig struct{
	Data struct {
		Fields map[string]FieldConfig `json:"fields"`
	} `json:"data"`
}

type FieldConfig struct {
	Path    string      `json:"path"`
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
}