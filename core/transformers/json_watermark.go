package transformers

import (
	"conecto/core"
	"context"

	"github.com/tidwall/gjson"
)

type JsonWatermark struct {
	Path string
}

func (t *JsonWatermark) Transform(ctx context.Context, batch []core.Event) ([]core.Event, error) {

	for i := range batch {
		ts := gjson.GetBytes(batch[i].Payload, t.Path).Time()
		batch[i].Timestamp = ts
	}

	return batch, nil
}