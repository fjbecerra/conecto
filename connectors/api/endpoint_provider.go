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
	metadataKeys []string
}

func NewDinamicEndpointProvider(endPointTemplate string, metadataKeys []string) *DinamicEndPointProvider{
	return &DinamicEndPointProvider{
		endPointTemplate: endPointTemplate,
		metadataKeys: metadataKeys,
	}
}

func (d *DinamicEndPointProvider) Apply(ctx context.Context, connection connections.Connection) string {
	args := make([]interface{}, len(d.metadataKeys))
	for i, key := range d.metadataKeys {
		args[i] = connection.Metadata[key]
	}
	return fmt.Sprintf(d.endPointTemplate, args...)
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

