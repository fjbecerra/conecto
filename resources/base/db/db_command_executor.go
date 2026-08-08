package db

import (
	"conecto/core/commands"
	"context"
	"database/sql"
	"errors"
)


type SQLExecutor struct {
    OpenTransaction func()*sql.Tx
}

func (e *SQLExecutor) Execute(
    ctx context.Context,
    commands []commands.Command,
) error {
    tx:= e.OpenTransaction()
    defer tx.Rollback()

    for _, command:= range commands {
        switch c := command.(type) {
            case *SQLCommand:
                 if _, err := tx.ExecContext(
                    ctx,
                    c.Query,
                    c.Values...,
                ); err != nil {
                    return err
                }
            default:
                return errors.New("unsupported command")
        }
    }
    return tx.Commit()
}

func (e *SQLExecutor) Close() error{
	return nil
}