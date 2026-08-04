package states

import (
	"conecto/core/statestores"
	"conecto/core/commands"
	"context"
	"fmt"
	"sync"
)


type MemoryStateStore struct{
	Store map[string]any
	mu    sync.RWMutex
}

func (s *MemoryStateStore) Load(context context.Context, ID string, name string) (statestores.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cursor, exist := s.Store[ID].(statestores.State)
	if(!exist){
		return statestores.State{}, nil
	}
	
	return cursor, nil
}

func (s *MemoryStateStore) Save(context context.Context, ID string, name string, state statestores.State) ([]commands.Command, error) {
	fmt.Println("SAVE CALLED:", state.Cursor)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Store[ID] = state
	return nil,nil
}

func (c *MemoryStateStore) Close() error{
	return nil
}