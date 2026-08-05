package api

import (
	"conecto/auth/connections"
)

type EndPointProvider interface {
	Apply(connection connections.Connection) string
}


