package rest

import (
	"context"
	"encoding/json"
	"testing"
)

func TestEmit5ElementsOverPaginating(t *testing.T) {	
	ctx := context.Background()
	mockClient := MockClient{
		calls: map[int]string {
				1:page1,
				2:page2,
				3:page3,
		},
	}
	paginationProvider := NewPaginationProvider(
		&mockClient,
		"../../../configs/facebook_ad_insight.json", 
	)

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


type MockClient struct {
	calls map[int]string
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

var pageCount = 0
func (m *MockClient) Fetch(ctx context.Context, url string) ([]byte, error) {
	pageCount++
	return extract(m.calls[pageCount])
	
}

func extract(body string) ([]byte, error) {
	return json.RawMessage([]byte(body)), nil
}

