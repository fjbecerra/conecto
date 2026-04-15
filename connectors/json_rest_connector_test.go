package connectors

import (
	"conecto/core"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestEmitTwoElements(t *testing.T) {	
	ctx := context.Background()
	json,_ := os.ReadFile("./testdata/ad_insight_action_product_id.json")
	mockClient := &http.Client{
		Transport: &MockRoundTripper{
			Body: string(json),
			StatusCode: 200,
		},
	}
	connector := JsonRestConnector{
    	client: mockClient,
		url: "http://fake-url",
		schema: core.LoadSchema("../pipelines/schemas/facebook_ad_insight.json"),
	}
	out, errCh := connector.Run(ctx)

	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
	
	var results []string

    for item := range out {
        results = append(results, string(item))
    }	

	if len(results) !=2 {
        t.Fatalf("expected 2 items, got %d", len(results))
    }
	
}

type MockRoundTripper struct {
    Body       string
    StatusCode int
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return &http.Response{
        StatusCode: m.StatusCode,
        Body: io.NopCloser(strings.NewReader(m.Body)),
        Header: make(http.Header),
    }, nil
}

