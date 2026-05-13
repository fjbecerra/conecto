package databases

type Adapter interface {
	Decode([]byte) (map[string]interface{}, error)
	BuildUpsertQuery(schema Schema,batchSize int) string

}