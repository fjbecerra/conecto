package shopify

import (
	"conecto/core"
	"conecto/core/statestores"
	"context"
	"time"

	"github.com/tidwall/gjson"
)

type JsonWatermark struct {
	Path string
}

func (t *JsonWatermark) Transform(ctx context.Context, batch *core.Batch) (*core.Batch, error) {
	var highestWatermark time.Time

	for i := range batch.Events {
		ts := gjson.GetBytes(batch.Events[i].Payload, t.Path).Time()
		if ts.After(highestWatermark) {
        	highestWatermark = ts
    	}
	}
	stateSyncState := statestores.SyncState{}
	stateSyncState["watermark"] = highestWatermark.Format(time.RFC3339Nano)

	return &core.Batch{
		Events: batch.Events,
		Cursor: batch.Cursor,
		IsLast: batch.IsLast,
		SyncState: stateSyncState,
	},nil
}