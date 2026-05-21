package core

import "context"

type Driver interface {
	Connection
	Executor
	BeginTx(context.Context) Transaction
}

type Connection interface {
	Close() error
}

type Transaction interface {
	Executor
	Commit() error
	Rollback() error
}

type Executor interface {
	Exec(context.Context, string, ...any) error
	Query(context.Context, string, ...any) (Rows, error)
	Prepare(context.Context, string) (Prepared, error)
}

// dependencies
type Prepared interface {
	Exec(context.Context, ...any) error
	Query(context.Context, ...any) (Rows, error)
}

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...any) error
}
