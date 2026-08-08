package api

import (
	"conecto/core"
)

type EndPointProvider interface {
	Apply(connection core.Connection) string
}


