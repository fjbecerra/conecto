package http

import "encoding/json"

type DataExtractor interface {

    Extract(body []byte,) ([]json.RawMessage,error)

}


