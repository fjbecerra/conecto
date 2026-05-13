package factories

import (
	"encoding/json"
	"os"
)

type ConnectorConfig struct{
	Type SourceType `json:"type"`
	RestConfig *RestConfig `json:"rest,omitempty"`
	MockedRestConfig *MockedRestConfig `json:"mocked_rest,omitempty"`
	Retry RetryConfig `json:"retry"`
}
type BaseRestConfig struct{
	Data struct {
		Path    string `json:"path"`
		IsArray bool   `json:"is_array"`
	}`json:"data"`
	Pagination struct {
		Type string `json:"type"`

		Request struct {
			Param string `json:"param"`
		} `json:"request"`

		Response struct {
			Next struct {
				Path string `json:"path"`
			} `json:"next"`
		} `json:"response"`
	} `json:"pagination"`

}
type RestConfig struct {
	BaseRestConfig
	BaseUrl string `json:"base_url"`
	Authentication struct{
		Type string `json:"type"`
		QueryToken struct {
			ParamName string `json:"param_name"`
		}`json:"query_token"`
	}`json:"authentication"`
}

type RetryConfig struct {
	MaxRetries int  `json:"max_retries"`
	BackoffMS  int  `json:"backoff_ms"`
	MaxBackoff int 	`json:"max_backoff"`
}

type MockedRestConfig struct{
	ResponsePaths []string `json:"response_paths"`
	BaseRestConfig
}

type TransformerConfig struct {
	Type TransformerType `json:"type"`
	ExtractorConfig *ExtractorConfig `json:"extractor,omitempty"`
}

type ExtractorConfig struct {	
	Fields string `json:"fields"`	
}

type SchemaConfig struct {
	DSN    string `json:"dsn"`
    Table 	string  `json:"table"`
	FieldsSpecs string `json:"fields_specs"`
	AutoCreate bool `json:"auto_create"`
}

type SinkConfig struct {
	Type  SinkType `json:"type"`
	SchemaConfig SchemaConfig `json:"schema"`
	BatchSize int `json:"batch_size"`
	Retry RetryConfig `json:"retry"`
}

type FieldsSpecs map[string]struct{
	Path    string      `json:"path"`
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
	Required bool		`json:"required"`
}

type RuntimeConfig struct {
	PipelineId 	string	`json:"pipeline_id"`
	StateStoreConfig StateStoreConfig `json:"state_store"`

}

type StateStoreConfig struct {
	Type	StateStoreType `json:"type"`
	Name 	string 			`json:"name"`
	AutoCreate bool `json:"auto_create"`
}

type DatabaseConfig struct {
	DSN string `json:"dsn"`
}


type ConfigPipeline struct {
	ConnectorConfig ConnectorConfig `json:"connector"`
	TransformersConfig []TransformerConfig `json:"transformers"`
	SinkConfig SinkConfig `json:"sink"`
	FieldsSpecsConfig map[string]FieldsSpecs `json:"fields_specs"`
	RuntimeConfig RuntimeConfig `json:"runtime"`
	DatabaseConfig DatabaseConfig `json:"database"`
}

func LoadConfigPipeline(path string) ConfigPipeline{
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    var config ConfigPipeline

    err = json.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }
		
    return config
}

type SourceType string
const (
	SourceRest 		 SourceType = "rest"
	SourceMockedRest SourceType = "mocked_rest"
)

type TransformerType string
const (
	TransformerExtractor TransformerType = "extractor"
)

type SinkType string
const (
	PostgresSink  SinkType = "postgres"
	MemorySink SinkType = "memory"
)

type StateStoreType string
const(
	MemoryStateStore StateStoreType = "memory"
	PostgresStateStore StateStoreType = "postgres"
)




