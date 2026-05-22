package stores

import (
	"context"
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
	SaveToken(context context.Context, ID string, record TokenRecord) error
	GetToken(context context.Context, ID string) (TokenRecord, error)
	Close()error
}