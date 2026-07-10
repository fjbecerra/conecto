package credentials

import (
	"context"
)

// type Metadata struct{
// 	ID string 
// 	provider string
// }

type CredentialService interface {
	Save(context context.Context, ID string, credential Credential) error
	Get(context context.Context, ID string) (Credential, error)
	Close() error
}


