package rest

import (
	"conecto/core"
	"fmt"
)

type RESTConnector struct {
	Provider *PaginationProvider
}

func (c *RESTConnector) FetchBatch(runtime core.Runtime, state core.Cursor) (core.Batch, error) {
	fmt.Println("SOURCE: sending event")
	var pc *PageCursor

	if state != nil {
		if v, ok := state["next"]; ok {
			pc = &PageCursor{Value: v}
		}
	}

	page, err := c.Provider.FetchPage(runtime, pc)
	if err != nil {
		return core.Batch{}, err
	}

	// convert raw rows → events
	events := make([]core.Event, len(page.Data))

	for i, row := range page.Data {
		events[i] = core.NewEvent(row)
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

func (c *RESTConnector) Close() error {	
    return c.Provider.Close()
}