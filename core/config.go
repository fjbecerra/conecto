package core

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseUrl string `json:"base_url"`
	Data struct {
		Path    string `json:"path"`
		IsArray bool   `json:"is_array"`
		FieldsConfig struct{
			Fields map[string]struct{
				Path    string      `json:"path"`
				Type    string      `json:"type"`
				Default interface{} `json:"default"`
			} `json:"fields"`
		}	
	}`json:"data"`

	Authentication struct{
		QueryToken struct {
			ParamName string `json:"param_name"`
		}`json:"query_token"`
	}`json:"authentication"`

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

func LoadConfig(path string) Config{
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    var config Config

    err = json.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }
    return config
}