package unittest

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"conecto/auth/oauth"
	"conecto/auth/oauth/state"
	"conecto/connectors"
	"conecto/connectors/shopify"
	"conecto/factories"
	"conecto/pipelines"
	"conecto/sync"
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var ctx = context.Background()
var router http.Handler
var credentialService credentials.CredentialService
var states string

func TestMain(m *testing.M) {

	stateSigner := state.NewHMACStateSigner(
			[]byte("z1Wbj51mq1GmIpmwAfLv9X5oSekOYEsC/9YXhOCuKjU="),
			10*time.Minute)
	states, _ = stateSigner.Sign("conn_123")
	v1, _ := base64.StdEncoding.DecodeString(
		"z1Wbj51mq1GmIpmwAfLv9X5oSekOYEsC/9YXhOCuKjU=",
	)
	keys:=map[string][]byte{
		"v1": v1,
	}
	keyManager:=credentials.NewStaticKeyManager(keys, "v1")

	credentialStore := credentials.NewMemoryStoreCredential(make(map[string]any))
	
	credentialService = credentials.NewAESGCMCredentialService(credentialStore, keyManager)
	router = setupRouterWithFakeConnector(ctx, credentialService, stateSigner)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestAuthorizeEndpoint(t *testing.T) {

	
	req := httptest.NewRequest(
		http.MethodGet,
		"/oauth/conn_123/authorize",
		nil,
	)

	rec := httptest.NewRecorder()


	router.ServeHTTP(
		rec,
		req,
	)


	require.Equal(
		t,
		http.StatusFound,
		rec.Code,
	)


	location :=
		rec.Header().Get("Location")


	require.Contains(
		t,
		location,
		"state=",
	)
}

func TestCallbackEndpoint(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/oauth/callback?code=test-code&state="+states,
		nil,
	)

	rec := httptest.NewRecorder()


	router.ServeHTTP(
		rec,
		req,
	)


	require.Equal(
		t,
		http.StatusFound,
		rec.Code,
	)


	require.Equal(
		t,
		"/connections",
		rec.Header().Get("Location"),
	)


	// Verify credential was stored
	credential, err := credentialService.Get(
		ctx,
		connections.Connection{
			ID: "conn_123",
		},
	)
	require.NoError(t, err)
	

	require.NoError(t, err)


	require.Equal(
		t,
		"test-token",
		credential.Data["access_token"],
	)
}

func setupRouterWithFakeConnector(context context.Context, credentialService credentials.CredentialService, stateSigner *state.HMACStateSigner) http.Handler {

	connectionStore := connections.NewMemoryStore()

	connectionStore.Save(context,connections.Connection{
		ID: "conn_123",
		Provider: "shopify",
		TenantID: "tenant_123",
		ExternalID: "external_id_123",
		Metadata: map[string]any{
			"shop": "xyz",
		},
		Status: "pending",
	})

	
	shopifyConnector := shopify.Connector{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Scopes:       []string{"read_orders"},
			HttpClient: &http.Client{
				Transport: &MockRoundTripper{
					Response: &http.Response{
						StatusCode: http.StatusOK,
						Header: http.Header{
							"X-Shopify-Access-Token": []string{"test-token"},
						},
					},
				},
			},
	}		
	registry := connectors.NewRegistry(shopifyConnector)

	queue := sync.NewQueue(10)
	config:= factories.LoadConfigPipeline("./testdata/orders_pipeline_flattened_data_to_memory.json")
	pipeline:= factories.BuildPipeline(config)
	pipeRegistry:= pipelines.NewRegistry()
	pipeRegistry.Register(pipeline)
	jobRepository:= sync.NewMemoryJobRepository()

	syncService := sync.NewSyncService(queue, pipeRegistry, connectionStore, jobRepository)


	service := oauth.NewService(
				connectionStore,
				credentialService,
				stateSigner,
				*registry, *syncService)


	handler := oauth.NewHandler(service)

	return oauth.NewRouter(handler)
}

type MockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}