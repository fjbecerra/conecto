package graphql

import (
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
		TokenProvider: &auths.HeaderTokenProvider{
			HeaderName: "any-header",
		},
		TokenStore: tokenStore,
	}

	builder := GraphQLRequestBuilder{
		Endpoint: "http://anyurl.com",
		Query: `query {
  orders(first: 10) {
    edges {
      cursor
      node {
        id
      }
    }

    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}`,
	}

	dataExtractor := GraphQLDataExtractor{
		Path: "data.orders.edges",
	}

	cursorExtractor := GraphQLCursorExtractor{
		HasMorePath: "data.orders.pageInfo.hasNextPage",
		CursorPath:  "data.orders.pageInfo.endCursor",
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

	if len(out.Events) != 1 {
		t.Fatalf("expected 1 items, got %d", len(out.Events))
	}

}

var page1 = `{
  "data": {
    "orders": {
      "edges": [
        {
          "cursor": "cursor-1",
          "node": {
            "id": "gid://shopify/Order/1"
          }
        }
      ],
      "pageInfo": {
        "hasNextPage": true,
        "endCursor": "cursor-1"
      }
    }
  }
}`
