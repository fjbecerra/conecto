package testutils

import (
	"context"
	"encoding/json"
)

type MockClient struct {
	Calls map[int]string
}

var pageCount = 0
func (m *MockClient) Fetch(ctx context.Context, url string) ([]byte, error) {
	pageCount++
	return extract(m.Calls[pageCount])
	
}

func extract(body string) ([]byte, error) {
	return json.RawMessage([]byte(body)), nil
}

