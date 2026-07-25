package api

import "net/http"

type HttpResponse struct {
	Body    []byte
	Headers http.Header
	Status  int
}
