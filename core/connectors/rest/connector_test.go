package rest

import (
	"conecto/core"
	"conecto/testutils"
	"context"
	"testing"
)

func TestEmit5ElementsOverPaginating(t *testing.T) {	
	ctx := context.Background()
	mockClient := testutils.MockHttpClient{
		Calls: map[int]string {
				0:page1,
		},
	}
  restClient:= RestClient{
    Client: &mockClient,
  }
	paginationProvider := PaginationProvider{
		RestClient : restClient,
    BaseUrl: "http://anyurl.com",
    DataPath: "data",
    ResponseNextPath: "paging.cursors.after",
		RequestParam: "after",
  }

	connector := RESTConnector {
		Provider: &paginationProvider,
	}

  runtime:= core.Runtime{
    Context:ctx,
  }
	out, _ := connector.FetchBatch(runtime, core.Cursor{})   

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
