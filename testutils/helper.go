package testutils

import (
	"io"
	"net/http"
	"os"
	"strings"
)

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

func MockHttpClient(jsonReponsePath string) *http.Client{
	json,_ := os.ReadFile(jsonReponsePath)
	return &http.Client{
		Transport: &MockRoundTripper{
			Body: string(json),
			StatusCode: 200,
		},
	}
}

