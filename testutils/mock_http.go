package testutils

import (
	"encoding/json"
	"net/http"
)

type MockHttpClient struct {
	Calls map[int]string

}

var pageCount = -1
func (m *MockHttpClient) Fetch(req *http.Request) ([]byte, error){
	pageCount++
	return extract(m.Calls[pageCount])
	
}

func (m *MockHttpClient) Close() error {
	return nil	
}

func extract(body string) ([]byte, error) {
	return json.RawMessage([]byte(body)), nil
}



