package extractors

import (
	"conecto/core"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var extractor = NewJsonExtractor(core.LoadSchema("../../pipelines/schemas/facebook_ad_insight.json"))
var ctx = context.Background()

func TestEmptyResponse(t *testing.T) {	
	in := make(chan json.RawMessage, 1)
	in <- []byte("")
	_, errCh := extractor.Extract(ctx, in)
	var errors []error
	for er := range errCh{
		errors = append(errors, er)
	}	
	for _, er := range errors{
		if er.Error() != "no json to parse" {
			t.Errorf("unexpected error: %v", er)
		}
	}
}

func TestExtractValuesFromFieldsOfFbInsightAd(t *testing.T) {
	in := make(chan json.RawMessage, 1)

    json,_ := os.ReadFile("./testdata/item_fb_insight_ad.json")	
    in <- json
    close(in) // 🔥 critical

    recordsCh, errCh := extractor.Extract(ctx, in)

    var records []core.Record

    for record := range recordsCh {
        fmt.Printf("RECEIVED: %+v\n", record)
        records = append(records, record)
    }

    for err := range errCh {
        if err != nil {
            t.Fatal(err)
        }
    }

    if len(records[0]) != 3 {
        t.Errorf("expected 3 records, received %d", len(records))
    }
	
}

