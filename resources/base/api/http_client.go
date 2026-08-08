package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IClient interface {
	Fetch(req *http.Request) (*HttpResponse, error)
	Close() error
}

type HttpClient struct {
	Client *http.Client
}

func (c *HttpClient) Fetch(req *http.Request) (*HttpResponse, error) {

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Body:    body,
		Headers: resp.Header,
		Status:  resp.StatusCode,
	}, nil
}

func (c *HttpClient) Post(url string, bodyRequest map[string]string) (*HttpResponse, error){
	jsonData, err := json.Marshal(bodyRequest)
	
	if err != nil {
		panic(err)
	}

	resp,err :=c.Client.Post(
			fmt.Sprintf(url),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
	defer resp.Body.Close()
	
	if(err!= nil){
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Body:    body,
		Headers: resp.Header,
		Status:  resp.StatusCode,
	}, nil

}

func (c *HttpClient) Close() error {
	c.Client.CloseIdleConnections()
	return nil
}

