package transformers

import "github.com/tidwall/gjson"

type JSONSelector interface {
	Select(payload []byte, path string) ([]byte, error)
}

type GJSONSelector struct{}

func (s *GJSONSelector) Select(payload []byte, path string) ([]byte, error) {

	result := gjson.GetBytes(payload, path)

	if !result.Exists() {
		return nil, nil
	}

	return []byte(result.Raw), nil
}