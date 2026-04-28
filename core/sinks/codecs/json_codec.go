package codecs

import "encoding/json"

type JSONCodec struct{}

func (c *JSONCodec) Decode(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	return m, err
}