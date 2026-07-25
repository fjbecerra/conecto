package rest

import (
	"conecto/connectors/api"
	"conecto/connectors/api/auths"
	"conecto/connectors/api/auths/stores"
	"conecto/core/statestores"
	"context"
	"encoding/base64"
	"testing"
	"time"
)

func TestEmit5ElementsOverPaginating(t *testing.T) {
	ctx := context.Background()
	v1, _ := base64.StdEncoding.DecodeString(
		"cQlswP0ZFMtD9dxWFOSUIF0ms9nJ5zj4ttxX5QfsmOU=",
	)
	keys := map[string][]byte{
		"v1": v1,
	}
	keyManager := auths.NewStaticKeyManager(keys, "v1")

	tokenStore := auths.NewADBTokenStore(stores.NewMemoryStoreToken(make(map[string]any)), keyManager)
	token := auths.Token{
		AccessToken:  "any-token",
		RefreshToken: "any-refresh-token",
		Expiry:       time.Now(),
	}
	error := tokenStore.Save(ctx, "id", token)
	if error != nil {
		t.Error(error.Error())
	}

	mockClient := api.MockHttpClient{
		Calls: map[int]string{
			0: page1,
		},
	}
	client := api.Client{
		Client: &mockClient,
		TokenProvider: &auths.QueryTokenProvider{
			ParamName: "token",
		},
		TokenStore: tokenStore,
	}

	builder := RestRequestBuilder{
		BaseURL:     "http://anyurl.com",
		CursorParam: "after",
		Method:      "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	dataExtractor := RestDataExtractor{
		Path: "data",
	}

	cursorExtractor := api.JSONCursorExtractor{
		Path: "paging.cursors.after",
	}

	paginationProvider := api.PaginationProvider{
		Client:  &client,
		Builder: &builder,
		Data:    &dataExtractor,
		Cursor:  &cursorExtractor,
	}

	connector := api.HttpConnector{
		Provider: &paginationProvider,
	}

	out, _ := connector.FetchBatch(ctx, statestores.Cursor{}, "id")

	if len(out.Events) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out.Events))
	}

}

var page1 = `{
  "data": [
    {"clicks": 1},
    {"clicks": 2}
  ],
  "paging": {
    "cursors": {
      "after": "cursor-1"
    }
  }
}`
