package transformers

import (
	"conecto/core"
	"context"
	"encoding/json"
)

type Fields map[string]struct{
	Path    string    
	Type    string
	Default interface{}
    Required bool
}

type Extractor struct {
    Selector JSONSelector
    Fields Fields
}

func (e *Extractor) Transform(ctx context.Context, batch []core.Event) ([]core.Event, error) {
    out := make([]core.Event, 0, len(batch))

	for _, ev := range batch {

		newObj := make(map[string]any)

		for field, fieldSpec := range e.Fields {

			val, err := e.Selector.Select(ev.Payload, fieldSpec.Path)
			if err != nil {
				return nil, err
			}

			if val != nil {
				newObj[field] = json.RawMessage(val)
			}
		}

		b, _ := json.Marshal(newObj)

		ev.Payload = b
		out = append(out, ev)
	}

	return out, nil
}
