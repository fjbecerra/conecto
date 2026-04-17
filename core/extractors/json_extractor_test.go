package extractors

import (
	"conecto/core"
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

var extractor = NewJsonExtractor("../../configs/facebook_ad_insight.json")
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
	var raw = `{
		"spend" : 1.5
		"clicks" : 1,
		"impressions: 2
	}`
	rawMessage := json.RawMessage([]byte(raw))
    in <- rawMessage
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

    if len(records) != 1 {
        t.Errorf("expected 1 records, received %d", len(records))
    }

	if len(records[0]) != 3 {
		t.Errorf("expected 3 fields, recieved %d", len(records[0]))
	}	
}

