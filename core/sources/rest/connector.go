package rest

import (
	"context"
	"encoding/json"
)

type Connector struct {
	Provider *PaginationProvider
}

func (c *Connector) Run(ctx context.Context) (<-chan json.RawMessage, <-chan error) {

	out := make(chan json.RawMessage, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		var cursor *Cursor

		for {
			page, err := c.Provider.FetchPage(ctx, cursor)
			if err != nil {
				errCh <- err
				return
			}

			// stream PAGE rows
			for _, row := range page.Data {
				select {
				case out <- row:
				case <-ctx.Done():
					return
				}
			}

			if !page.HasMore || page.NextCursor == nil {
				break
			}

			cursor = page.NextCursor
		}
	}()

	return out, errCh
}