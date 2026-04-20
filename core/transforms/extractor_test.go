package transforms

import (
	"context"
	"encoding/json"
	"testing"
)

var ctx = context.Background()

func TestEmptyResponse(t *testing.T) {	
	extractor := Extractor{}
	in := make(chan json.RawMessage, 1)
	in <- []byte("")
	_, error := extractor.Extract([]byte(""))
	
	if error.Error() != "no json to parse" {
			t.Errorf("unexpected error: %v", error.Error())
	}
	
}

func TestExtractValuesFromFieldsOfFbInsightAd(t *testing.T) {
	extractor := Extractor{
		Fields:Fields{
			"spend": {
				Path: "spend",
				Type: "float64",
				Default : 0,
			},
			"clicks": {
				Path: "clicks",
				Type: "int",
				Default : 0,
			},
			"impressions": {
				Path: "impressions",
				Type: "int",
				Default : 0,
			},
		},
	}
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

