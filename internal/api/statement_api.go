package api

import "context"

type Statement interface {
	Binder[Statement]
	Prepare(...context.Context) (Prepared, error)
	Raw() (string, []any, error)
	Run(...context.Context) error
}

type Prepared interface {
	Binder[Prepared]
	Run(...context.Context) error
}

// ---

type Binder[T any] interface {
	Bind(any) T
	Values(...any) T
	Value(id string, v any) T
}

type OnlyBinder interface {
	Bind(any) OnlyBinder
	Values(...any) OnlyBinder
	Value(id string, v any) OnlyBinder
}
