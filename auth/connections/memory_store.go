package connections

import (
	"context"
	"errors"
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
