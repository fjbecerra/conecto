package codecs

type Codec interface {
	Decode([]byte) (map[string]interface{}, error)
}
