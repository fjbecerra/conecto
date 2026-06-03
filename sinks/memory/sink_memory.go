package memory

import (
	"conecto/core"
    "conecto/core/commands"
	"conecto/core/codecs"
	"context"
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

func (m *SinkMemory) WriteBatch(context context.Context, ID string, batch [] core.Event) ([]commands.Command, error){
    fmt.Println("SINK: received event")
     if m.OnWrite != nil {
            m.OnWrite()
        }
    rows := make([]map[string]interface{}, 0, len(batch))
     m.mu.Lock()
    for _, ev := range batch {
       
        rec, err := m.Codec.Decode(ev.Payload)        
        if err != nil {
            return nil,err
        }
        //metadata
		for key, value := range ev.Meta{
			rec[key]=value
		}
        rec["__pipeline_id"]=ID
		rows = append(rows, rec)
       
    }
     m.Mstore = append(m.Mstore, rows...)   
     m.mu.Unlock()
    fmt.Println("SINK: exiting Run()")
    return nil,nil
}


