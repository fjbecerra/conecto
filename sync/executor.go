package sync

import (
	"conecto/auth/connections"
	"conecto/core/pipelines"
	"context"
)

type Executor struct {
	registry         pipelines.Registry
	connectionStore  connections.Store
}

func NewExecutor(
	registry pipelines.Registry,
	connectionStore connections.Store,
) *Executor {
	return &Executor{
		registry:        registry,
		connectionStore: connectionStore,
	}
}

func (e *Executor) Execute(
	ctx context.Context,
	job SyncJob,
) error {

	conn, err := e.connectionStore.Get(
		ctx,
		job.ConnectionID,
	)

	if err != nil {
		return err
	}

	p, err := e.registry.Get(
		job.PipelineID,
	)

	if err != nil {
		return err
	}

	for _, stream := range p.Streams {

		err := stream.Run(
			ctx,
			conn,
		)

		if err != nil {
			return err
		}
	}

	return nil
}