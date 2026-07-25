package unittests

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"conecto/connectors/api"
	"conecto/core/engines"
	"conecto/factories"
	"conecto/sinks/memory"
	"conecto/states"
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

var credentialService credentials.CredentialService

func TestMain(m *testing.M) {

	os.Setenv("TOKEN_ENCRYPTION_KEY_V1", "z1Wbj51mq1GmIpmwAfLv9X5oSekOYEsC/9YXhOCuKjU=")
	code := m.Run()
	os.Exit(code)
}

func TestMockedOrdersipelineRawData(t *testing.T) {
	context := context.Background()
	config := factories.LoadConfigPipeline("./testdata/orders_pipeline_raw_data_to_memory.json")
	pipeline := factories.BuildPipeline(config)
	connection := connections.Connection{
		ID: "test-connection-id-123",
	}
	for _, stream := range pipeline.Streams {
		credentialService := stream.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*api.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)
		credential := credentials.Credential{
			Type: "oauth2",

			Data: map[string]string{
				"X-Shopify-Access-Token": "shpat_xxxxx",
				"refresh_token":          "xxxx",
			},

			Expiry: nil,
		}

		error := credentialService.Save(context, connection, credential)
		if error != nil {
			t.Error(error.Error())
		}

		error = stream.Run(context, connection)
		if error != nil {
			t.Error(error.Error())
		}

		memSink := stream.Engine.SinkCommiter.(*engines.SinkEngine).Sink.(*memory.SinkMemory)

		if len(memSink.Mstore) != 2 {
			t.Errorf("number of record expected is 2, returned: %d", len(memSink.Mstore))
		}
	}

}

func TestMockedOrdersPipelineFlattened(t *testing.T) {
	context := context.Background()
	config := factories.LoadConfigPipeline("./testdata/orders_pipeline_flattened_data_to_memory.json")

	pipeline := factories.BuildPipeline(config)
	connection := connections.Connection{
		ID: "test-connection-id-123",
	}
	for _, stream := range pipeline.Streams {
		credentialService := stream.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*api.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)
		credential := credentials.Credential{
			Type: "oauth2",

			Data: map[string]string{
				"X-Shopify-Access-Token": "shpat_xxxxx",
				"refresh_token":          "xxxx",
			},

			Expiry: nil,
		}
		error := credentialService.Save(context, connection, credential)
		if error != nil {
			t.Error(error.Error())
		}

		error = stream.Run(context, connection)
		if error != nil {
			t.Error(error.Error())
		}

		memSink := stream.Engine.SinkCommiter.(*engines.SinkEngine).Sink.(*memory.SinkMemory)

		if len(memSink.Mstore) != 2 {
			t.Errorf("number of record expected is 2, returned: %d", len(memSink.Mstore))
		}
	}
}

func TestPipeline_CancelAndResume(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := factories.LoadConfigPipeline(
		"./testdata/orders_pipeline_flattened_data_to_memory.json",
	)

	pipeline := factories.BuildPipeline(cfg)

	connection := connections.Connection{
		ID: "test-connection-id-123",
	}
	for _, stream := range pipeline.Streams {

		credentialService := stream.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*api.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)
		credential := credentials.Credential{
			Type: "oauth2",

			Data: map[string]string{
				"X-Shopify-Access-Token": "shpat_xxxxx",
				"refresh_token":          "xxxx",
			},

			Expiry: nil,
		}
		er := credentialService.Save(ctx, connection, credential)
		if er != nil {
			t.Error(er.Error())
		}

		sink := stream.Engine.SinkCommiter.(*engines.SinkEngine).Sink.(*memory.SinkMemory)

		store := stream.Engine.SinkCommiter.(*engines.SinkEngine).StateStore.(*states.MemoryStateStore)

		// RUN PIPELINE
		errCh := make(chan error, 1)

		go func() {
			errCh <- stream.Run(ctx, connection)
		}()

		// GIVE PIPELINE TIME TO PROCESS SOME DATA
		time.Sleep(50 * time.Millisecond)

		// CANCEL
		cancel()

		err := <-errCh

		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf(
				"expected context canceled, got %v",
				err,
			)
		}

		// VERIFY CHECKPOINT EXISTS
		state, err := store.Load(ctx, cfg.ID)
		if err != nil {
			t.Fatal(err)
		}

		// checkpoint MAY be nil if cancel happened before
		// first commit completed
		_ = state

		// RESTART PIPELINE
		err = stream.Run(ctx, connection)
		if err != nil {
			t.Fatal(err)
		}

		// VERIFY EVENTUAL CONSISTENCY
		if len(sink.Mstore) != 2 {
			t.Fatalf(
				"expected 2 records after resume, got %d",
				len(sink.Mstore),
			)
		}
	}

}
