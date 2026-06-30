package http

import "net/http"

type HttpResponse struct {
	Body    []byte
	Headers http.Header
	Status  int
}