package driver

import "context"

type Driver interface {
	Connection
	Executor
}

type Connection interface {
	Close() error
}

type Executor interface {
	Exec(context.Context, string, ...any) error
	Query(context.Context, string, ...any) (Rows, error)
}

// dependencies

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...any) error
}
