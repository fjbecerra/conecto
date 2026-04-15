package core

type FieldConfig struct {
	Path     string      `json:"path"`
	Type     string      `json:"type"`
	Default  interface{} `json:"default"`
	Nullable bool        `json:"nullable"`
}

type ComputedConfig struct {
    Expr    string  `json:"expr"`
    Type    string  `json:"type"`
    Default interface{} `json:"default"`
}

type SchemaConfig struct {
	Path        string                      `json:"path"`
    Is_array    bool                        `json:"is_array"`
	Fields      map[string]FieldConfig      `json:"fields"`
    Computed    map[string]ComputedConfig   `json:"computed"`
}

