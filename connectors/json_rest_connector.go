package connectors

import (
	"conecto/core"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

type TokenProvider struct {
    token string
}

type PaginationAdapter struct {

}

type JsonRestConnector struct{
	client *http.Client
    url    string
    schema core.SchemaConfig
}

func NewJsonRestConnector(client *http.Client, url string, responseSchema core.SchemaConfig) core.Connector[json.RawMessage]{
    
    // url := fmt.Sprintf(`https://graph.facebook.com/v19.0/{page_id}/insights
    //         ?metric=page_fans,page_fan_adds_unique,page_fan_removes_unique
    //         ,page_impressions,page_impressions_unique
    //         ,page_post_engagements,page_actions_post_reactions_like_total
    //         ,page_actions_post_reactions_total,page_actions_post_comments_total
    //         ,page_actions_post_shares_total,page_total_actions
    //         ,page_storytellers
    //         ,page_video_views,page_video_view_time
    //         &period=day
    //         &since=YYYY-MM-DD
    //         &until=YYYY-MM-DD
    //         &access_token=%s`,accessToken)
    
    return &JsonRestConnector{
        client: client,
        url: url,
        schema: responseSchema,
    }
}

func (jrc *JsonRestConnector) Run(ctx context.Context) (<-chan json.RawMessage, <-chan error) {
	out := make(chan json.RawMessage, 100)
    errCh := make(chan error, 1)

    go func() {
        defer close(out)
        defer close(errCh)

        req, err := http.NewRequestWithContext(ctx, "GET", jrc.url, nil)
        if err != nil {
            errCh <- err
            return
        }

        resp, err := jrc.client.Do(req)
        if err != nil {
            errCh <- err
            return
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            errCh <- err
            return
        }

        // 🔥 gjson extraction
        res := gjson.GetBytes(body, jrc.schema.Path)

        if !res.IsArray() {
            errCh <- fmt.Errorf("'data' is not an array")
            return
        }

        for _, item := range res.Array() {
            select {
            case out <- json.RawMessage(item.Raw):
            case <-ctx.Done():
                return
            }
        }
    }()

    return out, errCh
}
