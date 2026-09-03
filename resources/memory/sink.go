package memory

import (
	"conecto/core"
	"conecto/core/codecs"
	"conecto/core/commands"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/statestores"
	"context"
	"fmt"
	"sync"
)

type internalMemorySink struct{
    mu   sync.Mutex
    Mstore []map[string]any
    codec codecs.Codec
    OnWrite func()
}

func (m *internalMemorySink) WriteBatch(context context.Context, ID string, batch [] core.Event) ([]commands.Command, error){
    fmt.Println("SINK: received event")
     if m.OnWrite != nil {
            m.OnWrite()
        }
    rows := make([]map[string]interface{}, 0, len(batch))
     m.mu.Lock()
    for _, ev := range batch {
       
        rec, err := m.codec.Decode(ev.Payload)        
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

type MemorySink struct {   
    store []map[string]any
    retryExecutor *retry.Executor
	stateStore statestores.StateStore	

}

func CreateMemorySink(memorySink MemorySink ,OnWrite func()) engines.SinkCommiter {	
     var mu sync.Mutex
    sink := &internalMemorySink{
        mu : mu,
        Mstore: memorySink.store,
        codec: &codecs.JSONCodec{},
        OnWrite: OnWrite,
    }
    sinker:= &engines.Sinker{
        Sink: sink,
        Executor: &MemoryExecutor{
            Store: &[]core.Event{},
        },
    }
    return &engines.SinkEngine{
		Sinker: *sinker,
		SinkRetry: *memorySink.retryExecutor,
		StateStore: memorySink.stateStore,
	}
}
