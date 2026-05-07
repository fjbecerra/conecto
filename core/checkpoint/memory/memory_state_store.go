package memory

import (
	"conecto/core"
	"context"
)


type MemoryStateStore struct{
	Store map[string]core.State
}

func (s *MemoryStateStore) Load(ctx context.Context, pipelineID string) (core.State, error) {
	cursor, exist := s.Store[pipelineID]
	if(!exist){
		return core.State{}, nil
	}
	
	return cursor, nil
}

func (s *MemoryStateStore) Save(ctx context.Context,pipelineID string,state core.State) error {
	s.Store[pipelineID] = state
	return nil
}
