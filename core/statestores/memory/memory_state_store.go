package memory

import (
	"conecto/core"
	"context"
	"fmt"
	"sync"
)


type MemoryStateStore struct{
	Store map[string]core.State
	mu    sync.RWMutex
}

func (s *MemoryStateStore) Load(ctx context.Context, pipelineID string) (core.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cursor, exist := s.Store[pipelineID]
	if(!exist){
		return core.State{}, nil
	}
	
	return cursor, nil
}

func (s *MemoryStateStore) Save(ctx context.Context,pipelineID string,state core.State) error {
	fmt.Println("SAVE CALLED:", state.Cursor)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Store[pipelineID] = state
	return nil
}
