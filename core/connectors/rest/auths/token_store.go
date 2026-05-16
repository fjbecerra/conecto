package auths

import (
	"conecto/core"
)

type TokenStore interface {
	Save(runtime core.Runtime, token Token) error
	Get(runtime core.Runtime) (Token, error)
	Close() error
}


