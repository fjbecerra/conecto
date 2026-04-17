package pipelines

import (
	"conecto/core"
	"conecto/core/extractors"
	"conecto/core/sources/rest"
	"conecto/sinks"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestPipeline(t *testing.T) {
	ctx :=context.Background()
	configPath :="../configs/facebook_ad_insight.json"
	httpClient := MockHttpClient("./testdata/ad_insight_action_product_id.json")
	client := rest.NewRestClient(&httpClient)
	paginationProvider := rest.NewPaginationProvider(
		client,
		configPath)

	connector := rest.Connector {
		Provider: &paginationProvider,
	}

	extractor := extractors.NewJsonExtractor(configPath)


	source := &core.MapSource[json.RawMessage, core.Record]{
		Upstream: &connector,
		MapFn: func(raw json.RawMessage) core.Record {
			record,_ := extractor.Extract(raw)
			return record
		},
	}

	memorySink := sinks.NewMemorySink[core.Record]()


	pipeline:= &Pipeline[core.Record]{
		source: source,
		sink:   memorySink,
	}
	pipeline.Run(ctx)
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
        Body: io.NopCloser(strings.NewReader(m.Body)),
        Header: make(http.Header),
    }, nil
}

func MockHttpClient(jsonReponsePath string) http.Client{
	json,_ := os.ReadFile(jsonReponsePath)
	return http.Client{
		Transport: &MockRoundTripper{
			Body: string(json),
			StatusCode: 200,
		},
	}
}
