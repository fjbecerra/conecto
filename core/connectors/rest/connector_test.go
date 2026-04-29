package rest

import (
	"conecto/core"
	"conecto/testutils"
	"context"
	"testing"
)

func TestEmit5ElementsOverPaginating(t *testing.T) {	
	ctx := context.Background()
	mockClient := testutils.MockClient{
		Calls: map[int]string {
				0:page1,
		},
	}
	paginationProvider := PaginationProvider{
		Client : &mockClient,
    BaseUrl: "http://anyurl.com",
    DataPath: "data",
    ResponseNextPath: "paging.cursors.after",
		RequestParam: "after",
  }

	connector := RESTConnector {
		Provider: &paginationProvider,
	}

	out, _ := connector.FetchBatch(ctx, core.Cursor{})   

	if len(out.Events) !=2 {
        t.Fatalf("expected 2 items, got %d", len(out.Events))
    }
	
}


var page1 = `{
  "data": [
    {"clicks": 1},
    {"clicks": 2}
  ],
  "paging": {
    "cursors": {
      "after": "cursor-1"
    }
  }
}`

var page2 = `{
  "data": [
    {"clicks": 3},
    {"clicks": 4}
  ],
  "paging": {
    "cursors": {
      "after": "cursor-2"
    }
  }
}`

var page3 = `{
  "data": [
    {"clicks": 5}
  ],
  "paging": {}
}`
