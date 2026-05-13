package statestores

import (
	"conecto/core"
	"conecto/core/commands"
	"fmt"
	"sync"
)


type MemoryStateStore struct{
	Store map[string]core.State
	mu    sync.RWMutex
}

func (s *MemoryStateStore) Load(runtime core.Runtime) (core.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cursor, exist := s.Store[runtime.PipelineId]
	if(!exist){
		return core.State{}, nil
	}
	
	return cursor, nil
}

func (s *MemoryStateStore) Save(runtime core.Runtime,state core.State) ([]commands.Command, error) {
	fmt.Println("SAVE CALLED:", state.Cursor)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Store[runtime.PipelineId] = state
	return nil,nil
}

func (c *MemoryStateStore) Close() error{
	return nil
}