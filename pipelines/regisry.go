package pipelines

import (
	"conecto/core/pipeline"
	"fmt"
	"sync"
)

type Registry interface {
	Register(p pipeline.Pipeline) error
	Get(id string) (pipeline.Pipeline, error)
}

type registry struct {
	pipelines map[string]pipeline.Pipeline
	mu        sync.RWMutex
}


func NewRegistry() Registry {
	return &registry{
		pipelines: make(map[string]pipeline.Pipeline),
	}
}


func (r *registry) Register(p pipeline.Pipeline) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pipelines[p.ID]; exists {
		return fmt.Errorf(
			"pipeline %s already registered",
			p.ID,
		)
	}


	r.pipelines[p.ID] = p

	return nil
}


func (r *registry) Get(id string) (pipeline.Pipeline, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()


	p, ok := r.pipelines[id]

	if !ok {
		return pipeline.Pipeline{}, fmt.Errorf(
			"pipeline %s not found",
			id,
		)
	}


	return p, nil
}