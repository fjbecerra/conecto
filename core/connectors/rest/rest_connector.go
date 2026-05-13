package rest

import (
	"conecto/core"
	"conecto/core/idempotency"
	"context"
	"fmt"
	"time"
)

type RESTConnector struct {
	Provider *PaginationProvider
	Generator idempotency.Generator
}

func (c *RESTConnector) FetchBatch(ctx context.Context, state core.Cursor) (core.Batch, error) {
	fmt.Println("SOURCE: sending event")
	var pc *PageCursor

	if state != nil {
		if v, ok := state["next"]; ok {
			pc = &PageCursor{Value: v}
		}
	}

	page, err := c.Provider.FetchPage(ctx, pc)
	if err != nil {
		return core.Batch{}, err
	}

	// convert raw rows → events
	events := make([]core.Event, len(page.Data))

	for i, row := range page.Data {
		events[i] = core.Event{
			Payload: row,
			Meta: core.EventMeta{
				"__event_id": c.Generator.Generate(row),
				"__ingested_at": time.Now(),
			},
		}
	}

	// compute next cursor
	var next core.Cursor
	if page.HasMore && page.NextCursor != nil {
		next = core.Cursor{
			"next": page.NextCursor.Value,
		}
	}

	return core.Batch{
		Events: events,
		Cursor: next,
	}, nil
}

func (c *RESTConnector) Open(ctx context.Context,state core.Cursor) error {
	return nil
}

func (c *RESTConnector) Close() error {	
    return c.Provider.Close()
}