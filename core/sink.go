package core

import (
	"conecto/core/commands"
	"context"
)


type Sink interface {
	WriteBatch(context context.Context, ID string, batch []Event) ([]commands.Command, error)
}