package db

import (
	"conecto/core/commands"
	"context"
	"database/sql"
)


type SQLExecutor struct {
    OpenTransaction func()*sql.Tx
}

func (e *SQLExecutor) Execute(
    ctx context.Context,
    command commands.Command,
) error {

    switch c := command.(type) {

    case *SQLCommand:
        tx:= e.OpenTransaction()
        defer tx.Rollback()
        tx.ExecContext(
            ctx,
            c.Query,
            c.Values...,
        )
        return tx.Commit()
    default:
        panic("unsupported command")
    }
}

func (e *SQLExecutor) Close() error{
	return nil
}