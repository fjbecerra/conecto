package clients

import (
	"conecto/resources/base/api"
	"net/http"
	"time"
)

type HttpClientConfig struct {
	Timeout int `json:"timeout_ms"`
	MaxIdleConnections int `json:"max_idle_connections"`
	MaxIdleConnectionsPerHost int `json:"max_idle_connections_per_host"`
	IdleConnectionTimeoutMs int `json:"idle_connection_timeout_ms"`
}

func CreateHttpClint(config *HttpClientConfig) *HttpClient{
	var client *HttpClient

	if(config !=nil){
		client = newHttpClient(config)
	}else{
		client =newHttpClient(nil)
	}
	return client
}



type HttpClient struct {
	client     *http.Client
}

func newHttpClient(config *HttpClientConfig) *HttpClient{
	var client http.Client
	if(config == nil){
		client = http.Client{
				Timeout:  time.Duration(config.Timeout) * time.Millisecond,
				Transport: &http.Transport{
					MaxIdleConns:        config.MaxIdleConnections,
					MaxIdleConnsPerHost: config.MaxIdleConnectionsPerHost,
					IdleConnTimeout:     time.Duration(config.IdleConnectionTimeoutMs) * time.Millisecond,
			},
		}
	}else{
		client = http.Client{}
	}

	
	return &HttpClient{
		client: &client,
	}
}

func (h *HttpClient) Get() *api.HttpClient {
	return &api.HttpClient{
		Client: h.client,
	}
}

func (h *HttpClient) Close() error {
	return h.Close()
}




