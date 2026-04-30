package core

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

type RDBSConfig struct {
	DBType RdbsType `json:"db_type"` // postgres, mysql, bigquery
	DSN    string `json:"dsn"`
    Table 	string  `json:"table"`
	Schema string `json:"schema"`
	Upsert string `json:"upsert"`
}

type SinkConfig struct {
	Type  SinkType `json:"type"`
	RDBSConfig RDBSConfig `json:"rdbs"`
	BatchSize int `json:"batch_size"`
	Retry RetryConfig `json:"retry"`
}

type FieldsConfig map[string]struct{
	Path    string      `json:"path"`
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
	Required bool		`json:"required"`
}

type AdditionalConfigRawMessage map[string]json.RawMessage

type AdditionalConfig struct{
	FieldsConfig map[string]FieldsConfig `json:"fields_configs"`
	BufferSize int `json:"event_buffer_size"`
}


type ConfigPipeline struct {
	ConnectorConfig ConnectorConfig `json:"connector"`
	TransformersConfig []TransformerConfig `json:"transformers"`
	SinkConfig SinkConfig `json:"sink"`
	AdditionalConfig AdditionalConfig `json:"additional_configs"`
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