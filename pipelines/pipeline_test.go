package pipelines

import (
	"conecto/core"
	"conecto/core/transforms"
	"conecto/core/sources/rest"
	"conecto/core/sinks"
	"conecto/testutils"
	"context"
	"encoding/json"
	"runtime"
	"testing"
)

func TestPipeline(t *testing.T) {
	ctx :=context.Background()
	config := core.LoadConfig("../configs/facebook_ad_insight.json")
	mockClient := testutils.MockClient{
		Calls: map[int]string {
				1:page1,
				2:page2,
		},
	}
	paginationProvider := rest.PaginationProvider{
		Client: &mockClient,
		BaseUrl: config.BaseUrl,
		DataPath: config.Data.Path,
		ResponseNextPath: config.Pagination.Response.Next.Path,
		RequestParam: config.Pagination.Request.Param,
	}

	connector := rest.Connector {
		Provider: &paginationProvider,
	}

	extractor := transforms.Extractor{
		Fields: transforms.Fields(config.Data.FieldsConfig.Fields),
	}


	source := &transforms.MapSource[json.RawMessage, core.Record]{
		Upstream: &connector,
		MapFn: func(raw json.RawMessage) core.Record {
			record,_ := extractor.Extract(raw)
			return record
		},
		Workers: runtime.NumCPU(),
	}

	memorySink := sinks.NewMemorySink[core.Record]()


	pipeline:= &Pipeline[core.Record]{
		source: source,
		sink:   memorySink,
	}
	pipeline.Run(ctx)
	data:= memorySink.Data()
	if len(data) != 3 {
		t.Errorf("shoud return %d records", 3)
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
    {"clicks": 5}
  ],
  "paging": {}
}`