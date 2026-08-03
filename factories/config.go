package factories

import (
	"encoding/json"
	"os"
)

type AuthorizeConfig struct {
	AuthorizeType AuthorizeType `json:"type"`
	Oauth Oauth `json:"oauth"`
}

type Oauth struct {
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scopes []string `json:"scopes"`
	AppUrl string `json:"app_url"`
}

type StreamConfig struct {
	ID  string `json:"id"`
	MockedRestConfig *MockedRestConfig `json:"mocked_rest,omitempty"`
	Query 	string `json:"query"`
	TransformersConfig []TransformerConfig `json:"transformers"`
	DestinationConfig DestinationConfig `json:"destination"`
	FieldsSpecsConfig map[string]FieldsSpecs `json:"fields_specs"`

}

type DestinationConfig struct {
	Name string `json:"name"`
	Keys string `json:"keys"`
	Schema *string `json:"schema,omitempty"`
}

type ApiConfig struct {
	Type ApiType `json:"type"`	
	EndpointConfig EndpointConfig `json:"endpoint"`
	RestConfig *RestConfig `json:"rest,omitempty"`
	GraphqlConfig *GraphqlConfig `json:"graphql,omitempty"`
	Source string `json:"source"`
}

type ConnectorConfig struct{
	Type ConnectorType `json:"type"`
	ApiConfig *ApiConfig `json:"api,omitempty"`
	Retry RetryConfig `json:"retry"`
	StreamsConfig []StreamConfig `json:"streams"`

}

type RestConfig struct{	
	DataConfig DataConfig `json:"data"`
    PaginationConfig PaginationConfig `json:"pagination"`
	AuthenticationConfig AuthenticationConfig `json:"authentication"`	
}

type GraphqlConfig struct {
	DataConfig DataConfig `json:"data"`
    PaginationConfig GraphqlPaginationConfig `json:"pagination"`
	AuthenticationConfig AuthenticationConfig `json:"authentication"`
}

type EndpointConfig struct {
	EndpointType EndpointType `json:"type"`
	Base string `json:"base"`
	MetadataKey string `json:"metadata_key"`
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



type DBConfig struct {
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

type SourcesConfig map[string]struct {
	DSN  string `json:"dsn"`
	Type SourcesType `json:"type"` 
}


type PipelineConfig struct {
	ID string `json:"id"`
	AuthorizeConfig AuthorizeConfig `json:"authorize"`
	ConnectorConfig ConnectorConfig `json:"connector"`
	SinkConfig SinkConfig `json:"sink"`
}

type SyncConfig struct{
	Buffer SyncBuffer `json:"buffer"`
	Scheduler SyncScheduler `json:"scheduler"`
	Retry RetryConfig `json:"retry"`
}

type SyncBuffer struct {
	BufferType BufferType `json:"type"`
	Size int `json:"size"`
}

type SyncScheduler struct {
	Duration string `json:"duration"`
}

type HttpServerConfig struct {
	Port int `json:"port"`
}

type ConectoConfig struct {
	SourcesConfig SourcesConfig `json:"sources"`
	SyncConfig SyncConfig `json:"sync"`
	PipelineRegistryConfig []string `json:"pipeline_registry"`
	DBConfig DBConfig `json:"db"`
	HttpServerConfig HttpServerConfig `json:"http_server"`
} 

type AppConfig struct {
	ConectoConfig ConectoConfig `json:"app"`
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

type SourcesType string
const(
	PostgresSource 	SourcesType = "postgres"
	MemorySource 	SourcesType = "memory"
	HttpSource SourcesType = "http"
)

type AuthenticationType string
const(
	Query AuthenticationType = "query"
	Header AuthenticationType = "header"
	Bearer AuthenticationType = "bearer"
)

type BufferType string
const(
	QueueType BufferType = "queue"
)

type EndpointType string
const(
	DinamicEndpointType EndpointType = "dinamic"
	StaticEndpointType EndpointType = "fixed"
)


type AuthorizeType string
const (
	OauthType AuthorizeType = "oauth"
)



