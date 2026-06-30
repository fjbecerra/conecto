package http

import (
	"conecto/core"
	"conecto/core/statestores"
	"context"
	"fmt"
)

type RESTConnector struct {
	Provider *PaginationProvider
}

func (c *RESTConnector) FetchBatch(context context.Context, state statestores.Cursor, ID string) (core.Batch, error) {
	fmt.Println("SOURCE: sending event")
	var pc *PageCursor

	if state != nil {
		if v, ok := state["next"]; ok {
			pc = &PageCursor{Value: v}
		}
	}

	page, err := c.Provider.FetchPage(context, pc, ID)
	if err != nil {
		return core.Batch{}, err
	}

	// convert raw rows → events
	events := make([]core.Event, len(page.Data))

	for i, row := range page.Data {
		events[i] = core.NewEvent(row)
	}

	// compute next cursor
	var next statestores.Cursor
	if page.HasMore && page.NextCursor != nil {
		next = statestores.Cursor{
			"next": page.NextCursor.Value,
		}
	}

	return core.Batch{
		Events: events,
		Cursor: next,
	}, nil
}

func (c *RESTConnector) Close() error {	
    return c.Provider.Close()
}