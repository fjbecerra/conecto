package api

import "conecto/core/statestores"

type SyncRequestProvider interface {
	Apply(syncState statestores.SyncState)map[string]any
}