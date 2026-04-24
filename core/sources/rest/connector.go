package rest

import (
	"context"
	"encoding/json"
	"fmt"
)

type Connector struct {
	Provider *PaginationProvider
}

func (c *Connector) Fetch(ctx context.Context) (<-chan json.RawMessage, <-chan error) {

	out := make(chan json.RawMessage)
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
		fmt.Println("source done, closing channel")
	}()

	return out, errCh
}