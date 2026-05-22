package factories

import (
	"encoding/json"
	"os"
)

type ConnectorConfig struct{
	Type ConnectorType `json:"type"`	
	Retry RetryConfig `json:"retry"`
	RestConfig *RestConfig `json:"rest,omitempty"`
	MockedRestConfig *MockedRestConfig `json:"mocked_rest,omitempty"`
}
type RestConfig struct{	
	BaseUrl string `json:"base_url"`
	DataConfig DataConfig `json:"data"`
    PaginationConfig PaginationConfig `json:"pagination"`
	TokenStoreConfig TokenStoreConfig `json:"token_store"`
	AuthenticationConfig AuthenticationConfig `json:"authentication"`	
}

type DataConfig struct {
	Path    string `json:"path"`
	IsArray bool   `json:"is_array"`
}

type PaginationConfig struct {
	Type string `json:"type"`

	Request struct {
		Param string `json:"param"`
	} `json:"request"`

	Response struct {
		Next struct {
			Path string `json:"path"`
		} `json:"next"`
	} `json:"response"`
}

type TokenStoreConfig struct {
	Type TokenStoreType `json:"type"`
	AutoCreate bool `json:"auto_create"`
	Name string `json:"name"`
	Source string `json:"source"`
}

type AuthenticationConfig struct{
	Type string `json:"type"`
	QueryToken struct {
		ParamName string `json:"param_name"`
	}`json:"query_token"`
}

type MockedRestConfig struct{
	ResponsePaths []string `json:"response_paths"`
}

type RetryConfig struct {
	MaxRetries int  `json:"max_retries"`
	BackoffMS  int  `json:"backoff_ms"`
	MaxBackoff int 	`json:"max_backoff"`
}


type TransformerConfig struct {
	Type TransformerType `json:"type"`
	ExtractorConfig *ExtractorConfig `json:"extractor,omitempty"`
}

type ExtractorConfig struct {	
	Fields string `json:"fields"`	
}

type SchemaConfig struct {
    Table 	string  `json:"table"`
	FieldsSpecs string `json:"fields_specs"`
	AutoCreate bool `json:"auto_create"`
}

type SinkConfig struct {
	Type  SinkType `json:"type"`
	SchemaConfig SchemaConfig `json:"schema"`
	BatchSize int `json:"batch_size"`
	Retry RetryConfig `json:"retry"`
	Source string `json:"source"`
}

type FieldsSpecs map[string]struct{
	Path    string      `json:"path"`
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
	Required bool		`json:"required"`
}

type RuntimeConfig struct {
	StateStoreConfig StateStoreConfig `json:"state_store"`
}

type StateStoreConfig struct {
	Type	StateStoreType `json:"type"`
	Name 	string 			`json:"name"`
	AutoCreate bool `json:"auto_create"`
	Source string `json:"source"`
}

type SourcesConfig map[string]struct {
	DSN  string `json:"dsn"`
	Type SourcesType `json:"type"` 
}


type ConfigPipeline struct {
	ID string `json:"pipeline_id"`
	ConnectorConfig ConnectorConfig `json:"connector"`
	TransformersConfig []TransformerConfig `json:"transformers"`
	SinkConfig SinkConfig `json:"sink"`
	FieldsSpecsConfig map[string]FieldsSpecs `json:"fields_specs"`
	RuntimeConfig RuntimeConfig `json:"runtime"`
	SourcesConfig SourcesConfig `json:"sources"`
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

type ConnectorType string
const (
	Rest 	   ConnectorType = "rest"
	MockedRest ConnectorType = "mocked_rest"
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

type TokenStoreType string
const(
	MemoryTokenStore TokenStoreType = "memory"
	PostgresTokenStore TokenStoreType = "postgres"
)

type SourcesType string
const(
	PostgresSource 	SourcesType = "postgres"
)




