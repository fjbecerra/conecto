package sync

import (
	"conecto/auth/connections"
	"conecto/pipelines"
	"context"
	"testing"
	"time"
    
	"github.com/stretchr/testify/require"
)

func TestScheduleDueSyncs(t *testing.T) {

    ctx := context.Background()

    buffer := NewQueue(1)

    connectionStore := connections.NewMemoryStore()

    jobRepo := NewMemoryJobRepository()

    registry := pipelines.NewRegistry()

    registry.Register(pipeline.Pipeline{
        ID: "shopify",
    })

    conn := connections.Connection{
        ID:         "conn_1",
        Provider:   "shopify",
        Status:     "connected",
        SyncStatus: "idle",
        NextSyncAt: time.Now().Add(-time.Minute),
    }

    connectionStore.Save(ctx, conn)

    service := &SyncService{
        buffer:           buffer,
        registry:        registry,
        connectionStore: connectionStore,
        jobRepository:   jobRepo,
    }

    err := service.ScheduleDueSyncs(ctx)
    require.NoError(t, err)

    select {
    case job := <-buffer.Consume():

        require.Equal(t, conn.ID, job.ConnectionID)
        require.Equal(t, "shopify", job.PipelineID)
        require.Equal(t, JobPending, job.Status)

    default:
        t.Fatal("expected one job in queue")
    }
}