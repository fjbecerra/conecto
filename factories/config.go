package factories

import (
	"encoding/json"
	"os"
)

type StreamConfig struct {
	Id  string `json:"stream_id"`
	MockedRestConfig *MockedRestConfig `json:"mocked_rest,omitempty"`
	TransformersConfig []TransformerConfig `json:"transformers"`
	DestinationConfig DestinationConfig `json:"destination"`
	FieldsSpecsConfig map[string]FieldsSpecs `json:"fields_specs"`

}

type DestinationConfig struct {
	Name string `json:"destination"`
	Keys string `json:"keys"`
	Schema *string `json:"schema,omitempty"`
}

type ApiConfig struct {
	Type ApiType `json:"type"`	
	RestConfig *RestConfig `json:"rest,omitempty"`
	GraphqlConfig *GraphqlConfig `json:"graphql,omitempty"`
}

type ConnectorConfig struct{
	Type ConnectorType `json:"type"`
	ApiConfig *ApiConfig `json:"api,omitempty"`
	Retry RetryConfig `json:"retry"`

}

type RestConfig struct{	
	BaseUrl string `json:"base_url"`
	DataConfig DataConfig `json:"data"`
    PaginationConfig PaginationConfig `json:"pagination"`
	TokenStoreConfig TokenStoreConfig `json:"token_store"`
	AuthenticationConfig AuthenticationConfig `json:"authentication"`	
}

type GraphqlConfig struct {
	BaseUrl string `json:"base_url"`
	Query 	string `json:"query"`
	DataConfig DataConfig `json:"data"`
    PaginationConfig GraphqlPaginationConfig `json:"pagination"`
	TokenStoreConfig TokenStoreConfig `json:"token_store"`
	AuthenticationConfig AuthenticationConfig `json:"authentication"`
}

type GraphqlPaginationConfig struct {
	HasMorePath string `json:"has_more_path"`
	CursorPath  string `json:"cursor_path"`
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
	Type AuthenticationType `json:"type"`
	ParamName string `json:"param_name"`	
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


type PipelineConfig struct {
	ID string `json:"id"`
	ConnectorConfig ConnectorConfig `json:"connector"`
	SinkConfig SinkConfig `json:"sink"`
	StreamsConfig []StreamConfig `json:"streams"`
}

type SyncConfig struct{}

type ConfigApp struct {
	RuntimeConfig RuntimeConfig `json:"runtime"`
	SourcesConfig SourcesConfig `json:"sources"`
	SyncConfig SyncConfig `json:"sync"`
	PipelinesConfig []string `json:"pipelines"`
}



func LoadConfig[T any](path string) (T, error) {
	var config T

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

type ConnectorType string 
const (
	Api 	ConnectorType = "api"
)

type ApiType string
const (
	Rest 	   ApiType = "rest"
	Graphql	   ApiType = "graphql"
	MockedRest ApiType = "mocked_rest"
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

type AuthenticationType string
const(
	Query AuthenticationType = "query"
	Header AuthenticationType = "header"
	Bearer AuthenticationType = "bearer"
)




