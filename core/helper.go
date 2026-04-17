package core

import (
	"encoding/json"
	"os"
)

func LoadConfig[T any](path string) T{
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    var config T

    err = json.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }
    return config
}