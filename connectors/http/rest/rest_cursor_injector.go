package rest

import (
	"conecto/connectors"
	"net/url"
)

type RestCursorInjector struct {
    Values url.Values
    Param string
}

func (r *RestCursorInjector) Inject(c * connectors.PageCursor) error {

    if c == nil {
        return nil
    }

    r.Values.Set(r.Param, c.Value)

    return nil
}