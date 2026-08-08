package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Stream struct {
	Name  string `json:"name"`
	Testing *Testing `json:"testing,omitempty"`
	Inbound Inbound `json:"inbound"`
	Outbound Outbound `json:"outbound"`
	FieldsSpecs FieldsSpecs `json:"fields_specs"`
}

type Streams struct {
	Streams []Stream `json:"streams"`
}

type Testing struct {
	MockedRest *MockedRest `json:"mocked_rest,omitempty"`
}

type Inbound struct {
	Resource string `json:"resource"`
	Config 	 ConfigBytes `json:"config"`
	
}

type Outbound struct {
	Resource string `json:"resource"`
	Config 	 ConfigBytes `json:"config"`
	Retry *Retry `json:"retry,omitempty"`
	BatchSize int `json:"batch_size,omitempty"`
}


type Store struct {
	Type string `json:"type"`
	Client ConfigBytes `json:"client"`
}

type MockedRest struct{
	ResponsePaths []string `json:"response_paths"`
}

type Retry struct {
	MaxRetries int  `json:"max_retries"`
	BackoffMS  int  `json:"backoff_ms"`
	MaxBackoff int 	`json:"max_backoff"`
}

type FieldsSpecs map[string]struct{
	Path    string      `json:"path"`
	Type    string      `json:"type"`
	Default interface{} `json:"default"`
	Required bool		`json:"required"`
}

type Resource struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Config *ConfigBytes `json:"config,omitempty"`
	Client  ConfigBytes `json:"client"`
	Retry *Retry `json:"retry,omitempty"`
}

type Sync struct{
	Buffer SyncBuffer `json:"buffer"`
	Scheduler SyncScheduler `json:"scheduler"`
	Retry Retry `json:"retry"`
}

type SyncBuffer struct {
	BufferType BufferType `json:"type"`
	Size int `json:"size"`
}

type SyncScheduler struct {
	Duration string `json:"duration"`
}

type HttpServer struct {
	Port int `json:"port"`
}

type Pipeline struct {
	Provider string `json:"provider"`
	Path string `json:"path"`
}

type Security struct {
    CredentialKey string `json:"credential_encryption_key"`
    StateSignerKey       string `json:"auth_state_signer_key"`
}

type Conecto struct {
	Resources []Resource `json:"resources"` 
	Sync Sync `json:"sync"`
	Pipelines []Pipeline `json:"pipelines"` 
	Store Store `json:"store"`
	HttpServer HttpServer `json:"http_server"`
	Security   Security   `json:"security"`
} 

type App struct {
	Conecto Conecto `json:"app"`
}


type Format string

const (
    FormatJSON Format = "json"
)

func Unmarshal[T any](data []byte, format Format) (T, error) {
    var config T

    var err error

    switch format {
    case FormatJSON:
        err = json.Unmarshal(data, &config)

    default:
        return config, fmt.Errorf(
            "unsupported config format: %s",
            format,
        )
    }

    if err != nil {
        return config, err
    }

    return config, nil
}

func formatFromPath(path string) (Format, error) {
    switch filepath.Ext(path) {
    case ".json":
        return FormatJSON, nil
    default:
        return "", fmt.Errorf(
            "unsupported config format: %s",
            filepath.Ext(path),
        )
    }
}

func loadConfig[T any](path string) ([]byte, Format, error) {
	format, err := formatFromPath(path)
    if err != nil {
        return nil, format, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        return nil, format, err
    }
	return data, format, err
}

func LoadPipelineConfig[T any](path string) (T, error) {
	data, format, err := loadConfig[T](path)
	if err != nil {
		var zero T
        return zero, err
    }
	 return Unmarshal[T](data, format)
}

func LoadConectoConfig[T any](path string) (T, error) {
    data, format, err := loadConfig[T](path)
	if err != nil {
		var zero T
        return zero, err
    }

	data = []byte(os.ExpandEnv(string(data)))
    if strings.Contains(string(data), "${") {
		var zero T
        return zero, fmt.Errorf(
            "config contains unresolved environment variables",
        )
    }

    return Unmarshal[T](data, format)
}

type BufferType string
const(
	QueueType BufferType = "queue"
)