package api

import (
	"conecto/auth/connections"
	"context"
	"fmt"
)

type EndPointProvider interface {
	Apply(ctx context.Context, connection connections.Connection) string
}

type DinamicEndPointProvider struct {
	endPointTemplate string
	metadataKey string
}

func NewDinamicEndpointProvider(endPointTemplate string, metadataKey string) *DinamicEndPointProvider{
	return &DinamicEndPointProvider{
		endPointTemplate: endPointTemplate,
		metadataKey: metadataKey,
	}
}

func (d *DinamicEndPointProvider) Apply(ctx context.Context, connection connections.Connection) string {
	return fmt.Sprintf(d.endPointTemplate, connection.Metadata[d.metadataKey])
}

type StaticEndpointProvider struct {
	endPoint string
}

func NewStaticEndpointProvider(endPoint string) *StaticEndpointProvider{
	return &StaticEndpointProvider{
		endPoint: endPoint,
	}
}

func (d *StaticEndpointProvider) Apply(ctx context.Context, connection connections.Connection) string {
	return d.endPoint
}

