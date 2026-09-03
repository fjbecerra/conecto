package api

type IncrementalSyncProvider interface {
	Apply(watermark *string) string
}