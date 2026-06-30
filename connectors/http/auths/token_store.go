package auths

import "context"

type Metadata struct{
	ID string 
	provider string
}

type TokenStore interface {
	Save(context context.Context, ID string, token Token) error
	Get(context context.Context, ID string) (Token, error)
	Close() error
}


