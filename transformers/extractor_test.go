package transformers

import (
	"conecto/core"
	"context"
	"testing"
)

var ctx = context.Background()

func TestEmptyResponse(t *testing.T) {
	selector := GJSONSelector{}
	fields := Fields{}
	extractor := Extractor{Selector: &selector, Fields: fields}
	_, error := extractor.Transform(ctx, []core.Event{})

	if error.Error() != "no batch to process." {
		t.Errorf("unexpected error: %v", error.Error())
	}

}

func TestExtractValuesFromFieldsOfFbInsightAd(t *testing.T) {
	extractor := Extractor{
		Fields: Fields{
			"spend": {
				Path:    "spend",
				Type:    "float64",
				Default: 0,
			},
			"clicks": {
				Path:    "clicks",
				Type:    "int",
				Default: 0,
			},
			"impressions": {
				Path:    "impressions",
				Type:    "int",
				Default: 0,
			},
		},
		Selector: &GJSONSelector{},
	}
	var raw = `{
		"spend" : 1.5
		"clicks" : 1,
		"impressions: 2
	}`
	events := []core.Event{}
	events = append(events, core.Event{
		Payload: []byte(raw),
	})

	records, _ := extractor.Transform(ctx, events)

	if len(records) != 1 {
		t.Errorf("expected 3 records, received %d", len(records))
	}

}

func TestEmptyFields(t *testing.T) {
	extractor := Extractor{
		Fields:   nil,
		Selector: &GJSONSelector{},
	}
	var raw = `{
		"spend" : 1.5
		"clicks" : 1,
		"impressions: 2
	}`
	events := []core.Event{}
	events = append(events, core.Event{
		Payload: []byte(raw),
	})

	records, error := extractor.Transform(ctx, events)

	if len(records) != 0 {
		t.Errorf("expected 0 records, received %d", len(records))
	}

	if error.Error() != "no fields specs found." {
		t.Errorf("unexpected error: %v", error.Error())
	}

}
