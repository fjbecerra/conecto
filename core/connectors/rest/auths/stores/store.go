package stores

import (
	"conecto/core"
	"errors"
	"time"
)

var ErrTokenNotFound = errors.New("token not found")

type TokenRecord struct {
	Ciphertext []byte
	Nonce      []byte
	KeyVersion string
	ExpiresAt  time.Time

}

type Store interface {
	SaveToken(runtime core.Runtime, record TokenRecord) error
	GetToken(runtime core.Runtime) (TokenRecord, error)
	Close()error
}