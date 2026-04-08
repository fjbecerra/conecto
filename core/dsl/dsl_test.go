package transformer

import (
	"testing"
)

func TestExtractValues(t *testing.T) {
	
    schema := LoadSchema()

    input := LoadResponse("ad_insight_action_product_id.json")

	

    output := Extract(input, schema)

    if output[0]["spend"] != 14.35 {
		t.Errorf("expected 14.35, got %d", output[0]["spend"])
	}

	if output[0]["clicks"] != 9 {
		t.Errorf("expected 9, got %d", output[0]["clicks"])
	}

	if output[0]["impressions"] != 185 {
		t.Errorf("expected 185, got %d", output[0]["impressions"])
	}

	if output[1]["spend"] != 8.53 {
		t.Errorf("expected 8.53, got %d", output[1]["spend"])
	}

	if output[1]["clicks"] != 7 {
		t.Errorf("expected 4, got %d", output[1]["clicks"])
	}

	if output[1]["impressions"] != 171 {
		t.Errorf("expected 185, got %d", output[1]["impressions"])
	}

}

func TestEvalExpressions(t *testing.T) {
	schema := LoadSchema()

    input := LoadResponse("ad_insight_action_product_id.json")	

    output := Extract(input, schema)

	
	if output[0]["ctr"] != 4.864864864864865 {
		t.Errorf("expected 4.864864864864865, got %d", output[0]["ctr"])
	}

	
	if output[1]["ctr"] != 4.093567251461988 {
		t.Errorf("expected 4.093567251461988, got %d", output[1]["ctr"])
	}

}