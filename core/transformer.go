package core

import (
	"context"
)

type Transformer interface {
	Transform(ctx context.Context, batch []Event) ([]Event, error)
}

type Chain struct {
	Steps []Transformer
}

func (c *Chain) Transform(ctx context.Context,batch []Event,) ([]Event, error) {

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