package pipelines

import (
	"conecto/connectors"
	"conecto/core"
	"conecto/core/extractors"
	"conecto/sinks"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestPipeline(t *testing.T) {
	json,_ := os.ReadFile("./testdata/ad_insight_action_product_id.json")
	client := &http.Client{
    	Transport: &MockRoundTripper{
			Body: string(json),
			StatusCode: 200,
    	},
	}
	schema:= core.LoadSchema("./schemas/facebook_ad_insight.json")
	memorySink:=sinks.NewMemorySink()
	pipeline := &JsonPipeline{
		connector: connectors.NewJsonRestConnector(
			client,
			"url",
			schema,
			
		),
		extractor: extractors.NewJsonExtractor(schema),
		sink:   memorySink,
	}
	pipeline.Run(context.Background())
	data:= memorySink.Data()
	if len(data) != 2 {
		t.Errorf("shoud return %d records", 2)
	}

}

type MockRoundTripper struct {
    Body       string
    StatusCode int
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return &http.Response{
        StatusCode: m.StatusCode,
        Body:       io.NopCloser(strings.NewReader(m.Body)),
        Header:     make(http.Header),
    }, nil
}