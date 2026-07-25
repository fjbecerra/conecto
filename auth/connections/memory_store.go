package connections

import (
	"context"
	"errors"
	"time"
)

var ErrConnectionNotFound = errors.New("connection not found")

type MemoryStore struct {
	store map[string]Connection
}


func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		store: make(map[string]Connection),
	}
}

func (s *MemoryStore) Get(ctx context.Context,id string) (Connection,error){

	connection, ok := s.store[id]

	if !ok {
		return Connection{}, ErrConnectionNotFound
	}
	return connection, nil
}

func (s *MemoryStore) Save(ctx context.Context,connection Connection) error {

	s.store[connection.ID] = connection
	return nil
}

func (s *MemoryStore) UpdateStatus(ctx context.Context,id string,status string) error {

	connection, ok := s.store[id]
	if !ok {
		return ErrConnectionNotFound
	}
	connection.Status = status
	s.store[id] = connection
	return nil
}

func (s *MemoryStore) ClaimDueConnections(ctx context.Context) ([]Connection, error) {

	now := time.Now()

	var result []Connection


	for _, conn := range s.store {

		if conn.Status != "connected" {
			continue
		}


		if conn.SyncStatus != "idle" {
			continue
		}


		if conn.NextSyncAt.After(now) {
			continue
		}

		conn.SyncStatus = "queued"

		result = append(result,conn)
	}


	return result, nil
}

func (s *MemoryStore) MarkRunning(ctx context.Context, id string) error {
	conn, ok := s.store[id]
	if !ok {
		return ErrConnectionNotFound
	}
	conn.SyncStatus = "running"
	return nil
}

func (s *MemoryStore) MarkCompleted(ctx context.Context, id string, next time.Time) error {
	conn, ok := s.store[id]

	if !ok {
		return ErrConnectionNotFound
	}
	now := time.Now()
	conn.SyncStatus = "idle"
	conn.LastSyncAt = &now
	conn.NextSyncAt = next
	return nil
}

