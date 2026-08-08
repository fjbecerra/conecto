package db

type Adapter interface {
	Decode([]byte) (map[string]interface{}, error)
	BuildQuery(schema Schema,batchSize int) string

}