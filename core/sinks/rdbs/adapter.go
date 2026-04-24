package rdbs

type Adapter interface {
	BuildUpsertQuery(schema Schema, upsert Upsert ,batchSize int) string

}