package entity

import "context"

type Driver interface {
	Exec(context.Context, string, ...any) error
	Query(context.Context, string, ...any) (Rows, error)
	// Begin(context.Context) (Tx, error)
}

// --- Driver dependencies

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(...any) error
}

// type Tx interface{}
