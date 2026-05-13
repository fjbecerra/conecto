package sinks

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/sinks/codecs"
	"fmt"
	"sync"
)

type SinkMemory struct {
    mu   sync.Mutex
    Mstore []map[string]interface{}
    Codec codecs.Codec
    OnWrite func()
}

func NewMemorySink(mstore[]map[string]interface{}, codec codecs.Codec) *SinkMemory{
    var mu sync.Mutex
    return &SinkMemory{
        mu : mu,
        Mstore: mstore,
        Codec: codec,
    }
}

func (m *SinkMemory) WriteBatch(runtime core.Runtime, batch [] core.Event) ([]commands.Command, error){
    fmt.Println("SINK: received event")
     if m.OnWrite != nil {
            m.OnWrite()
        }
    for _, event := range batch {
        m.mu.Lock()
        record, error := m.Codec.Decode(event.Payload)
        
        if error != nil {
            return nil,error
        }
        
        m.Mstore = append(m.Mstore, record)
       
        m.mu.Unlock()
    }
    fmt.Println("SINK: exiting Run()")
    return nil,nil
}


