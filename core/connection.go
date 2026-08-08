package core

import "time"

type ConnectionStatus string
const (
	StatusConnected ConnectionStatus = "connected"
	StatusUnconnected ConnectionStatus = "Unconnected"
	SyncIdle ConnectionStatus   = "idle"
	SyncQueued ConnectionStatus = "queued"
	SyncRunning ConnectionStatus = "running"
)

type Provider string
type Metadata map[string]any

type Connection struct {
	ID string
	TenantID string
	Provider Provider
	ResourceName string
	ExternalID string //The identifier used by the external system.
	Metadata Metadata //This lets each connector persist small pieces of configuration without changing the schema.
	Status ConnectionStatus //Lifecycle of the connection.

	SyncStatus string
	NextSyncAt time.Time
	LastSyncAt *time.Time
}

