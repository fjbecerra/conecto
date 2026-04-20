package rest

import (
	"context"
	"conecto/testutils"
	"testing"
)

func TestEmit5ElementsOverPaginating(t *testing.T) {	
	ctx := context.Background()
	mockClient := testutils.MockClient{
		Calls: map[int]string {
				1:page1,
				2:page2,
				3:page3,
		},
	}
	paginationProvider := PaginationProvider{
		Client : &mockClient,
    BaseUrl: "http://anyurl.com",
    DataPath: "data",
    ResponseNextPath: "paging.cursors.after",
		RequestParam: "after",
  }

	connector := Connector {
		Provider: &paginationProvider,
	}

	out, errCh := connector.Fetch(ctx)

	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
	
	var results []string

    for item := range out {
        results = append(results, string(item))
    }	

	if len(results) !=5 {
        t.Fatalf("expected 5 items, got %d", len(results))
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
