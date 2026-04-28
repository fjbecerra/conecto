package transformers

import (
	"conecto/core"
	"context"
)

type Chain struct {
	Steps []Transformer
}

func (c *Chain) Transform(
	ctx context.Context,
	batch []core.Event,
) ([]core.Event, error) {

	var err error
	current := batch

	for _, t := range c.Steps {
		current, err = t.Transform(ctx, current)
		if err != nil {
			return nil, err
		}
	}

	return current, nil
}