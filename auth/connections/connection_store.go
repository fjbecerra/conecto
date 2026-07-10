package connections

import "context"

type Store interface {

	Get(ctx context.Context,id string) (Connection,error)

	Save(ctx context.Context, connection Connection) error

	UpdateStatus(ctx context.Context,id string,status string) error
}