package commands

import "context"

type CommandExecutor interface{
	Execute(ctx context.Context,commands []Command) error
	Close() error
}

