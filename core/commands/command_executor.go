package commands

import "context"

type CommandExecutor interface{
	Execute(ctx context.Context,command Command) error
	Close() error
}

