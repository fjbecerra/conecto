package extractors

import (
	"context"
	"encoding/json"
	"testing"
)

var extractor = NewJsonExtractor("../../configs/facebook_ad_insight.json")
var ctx = context.Background()

func TestEmptyResponse(t *testing.T) {	
	in := make(chan json.RawMessage, 1)
	in <- []byte("")
	_, error := extractor.Extract([]byte(""))
	
	if error.Error() != "no json to parse" {
			t.Errorf("unexpected error: %v", error.Error())
	}
	
}

func TestExtractValuesFromFieldsOfFbInsightAd(t *testing.T) {
	var raw = `{
		"spend" : 1.5
		"clicks" : 1,
		"impressions: 2
	}`
	in := json.RawMessage([]byte(raw))    

    records, _ := extractor.Extract(in)

    if len(records) != 3 {
        t.Errorf("expected 1 records, received %d", len(records))
    }

	
}

