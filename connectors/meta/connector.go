// package meta

// import (
// 	"conecto/connectors"
// 	"conecto/connectors/credentials"
// 	"context"
// 	"net/url"
// 	"time"
// )


// type Connector struct {
// 	AppID string
// 	AppSecret string
// }

// func (c Connector) Name() string {

// 	return "meta"
// }

// func (c Connector) AuthorizeURL(ctx context.Context, connection connectors.Connection, state string)(string,error){
// 	params := url.Values{}

// 	params.Set(
// 		"client_id",
// 		c.AppID,
// 	)

// 	params.Set(
// 		"redirect_uri",
// 		"https://myapp.com/oauth/callback",
// 	)

// 	params.Set(
// 		"state",
// 		state,
// 	)

// 	params.Set(
// 		"scope",
// 		"ads_management",
// 	)

// 	return "https://www.facebook.com/vXX.X/dialog/oauth?" +
// 		params.Encode(),
// 		nil
// }

// func (c Connector) ExchangeCode(
// 	ctx context.Context,
// 	connection connectors.Connection,
// 	code string,
// )(credentials.Credential,error){


// 	// call Meta token endpoint


// 	return credentials.Credential{

// 		Type:"oauth2",

// 		Data:map[string]string{

// 			"access_token":
// 			"EAABxxx",

// 		},


// 		Expiry:
// 		time.Now().
// 			Add(
// 				60*24*time.Hour,
// 			),

// 	},nil
// }

// func (c Connector) Refresh(
// 	ctx context.Context,
// 	connection connectors.Connection,
// 	old credentials.Credential,

// )(
// credentials.Credential,error,
// ){


// 	return c.exchangeLongLivedToken(
// 		ctx,
// 		old.Data["access_token"],
// 	)
// }