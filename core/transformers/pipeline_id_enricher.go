package transformers

import (
	"conecto/core"
	"context"
	"encoding/json"
	"errors"
)


type PipelineIdEnricher struct {
    PipelineId string
}

func (e *PipelineIdEnricher) Transform(ctx context.Context, batch []core.Event) ([]core.Event, error) {
    out := make([]core.Event, 0, len(batch))

	if len(batch) == 0{
		return nil, errors.New("no batch to process.")
	}

	if(len(e.PipelineId) == 0) {
		return nil, errors.New("no pipeline id.")
	}

	for _, ev := range batch {		

		var obj map[string]any

		err := json.Unmarshal(ev.Payload, &obj)
		if err != nil {
			return nil, err
		}

		obj[core.PipelineId] = e.PipelineId
		b, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}

		ev.Payload = b
		out = append(out, ev)
	}

	return out, nil
}
