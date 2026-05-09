package factories

import (
	"conecto/core"
	"conecto/core/statestores"
	"conecto/core/statestores/memory"
)

type StateStore struct{
	StateStoreConfig StateStoreConfig
}

func NewStateStore(stateStoreConfig StateStoreConfig) *StateStore {
	return &StateStore {
		StateStoreConfig: stateStoreConfig,
	}
}

func (c *StateStore) Build() statestores.StateStore {
	if c.StateStoreConfig.Type == "" {
		c.StateStoreConfig.Type = MemoryStateStore
	}
	var stateStore statestores.StateStore
	switch c.StateStoreConfig.Type {
		case MemoryStateStore:
			stateStore = &memory.MemoryStateStore{
				Store :  map[string]core.State{},				
			}
		default: panic("Unkown state store type")
	}
	return stateStore
}
