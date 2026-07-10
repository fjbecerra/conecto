package state

import "time"

type AuthorizationState struct {
    ConnectionID string    `json:"connection_id"`
    Provider     string    `json:"provider"`
    ExpiresAt    time.Time `json:"expires_at"`
}

type StateSigner interface {

    Sign(connectionID string) (string, error)

    Verify(state string) (string, error)
}

