package connections

import (
	"conecto/core"
	"context"
	"errors"
	"time"
)

var ErrConnectionNotFound = errors.New("connection not found")

type MemoryStore struct {
	store map[string]any
}


func NewMemoryConnectionStore(store map[string]any) *MemoryStore {

	return &MemoryStore{
		store: store,
	}
}

func (s *MemoryStore) Get(ctx context.Context,id string) (core.Connection,error){

	connection, ok := s.store[id].(core.Connection)

	if !ok {
		return core.Connection{}, ErrConnectionNotFound
	}
	return connection, nil
}

func (s *MemoryStore) Save(ctx context.Context,connection core.Connection) error {

	s.store[connection.ID] = connection
	return nil
}

func (s *MemoryStore) UpdateStatus(ctx context.Context,id string,status core.ConnectionStatus) error {

	connection, ok := s.store[id].(core.Connection)
	if !ok {
		return ErrConnectionNotFound
	}
	connection.Status = status
	s.store[id] = connection
	return nil
}

func (s *MemoryStore) ClaimDueConnections(ctx context.Context) ([]core.Connection, error) {

	now := time.Now()

	var result []core.Connection


	for _, conn := range s.store {

		connection, _ := conn.(core.Connection)

		if connection.Status != "connected" {
			continue
		}


		if connection.SyncStatus != "idle" {
			continue
		}


		if connection.NextSyncAt.After(now) {
			continue
		}

		connection.SyncStatus = "queued"

		result = append(result,connection)
	}


	return result, nil
}

func (s *MemoryStore) MarkRunning(ctx context.Context, id string) error {
	conn, ok := s.store[id].(core.Connection)
	if !ok {
		return ErrConnectionNotFound
	}
	conn.SyncStatus = "running"
	return nil
}

func (s *MemoryStore) MarkCompleted(ctx context.Context, id string, next time.Time) error {
	conn, ok := s.store[id].(core.Connection)

	if !ok {
		return ErrConnectionNotFound
	}
	now := time.Now()
	conn.SyncStatus = "idle"
	conn.LastSyncAt = &now
	conn.NextSyncAt = next
	return nil
}

