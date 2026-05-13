package commands

import (
	"context"
	"database/sql"
)

type SQLExecutor struct {
    Tx *sql.Tx
}

func (e *SQLExecutor) Execute(
    ctx context.Context,
    command Command,
) error {

    switch c := command.(type) {

    case *SQLCommand:
        _, err := e.Tx.ExecContext(
            ctx,
            c.Query,
            c.Values...,
        )
        return err
    default:
        panic("unsupported command")
    }
}

func (e *SQLExecutor) Close() error{
	return nil
}