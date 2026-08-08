package api

import (
	"conecto/core"
	"conecto/core/statestores"
	"context"
	"fmt"
)

type HttpConnector struct {
	Provider *PaginationProvider
}

func (c *HttpConnector) FetchBatch(context context.Context, state statestores.Cursor, connection core.Connection) (core.Batch, error) {
	fmt.Println("SOURCE: sending event")
	var pc *PageCursor

	if state != nil {
		if v, ok := state["next"]; ok {
			pc = &PageCursor{Value: v}
		}
	}

	page, err := c.Provider.FetchPage(context, pc, connection)
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
