package rdbs

type Adapter interface {
	Decode([]byte) (map[string]interface{}, error)
	BuildUpsertQuery(schema Schema, upsert Upsert ,batchSize int) string

}