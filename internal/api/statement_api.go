package api

import "context"

type Statement interface {
	Runner
	Binder[Statement]
	Prepare(...context.Context) (Prepared, error)
}

type Prepared interface {
	Binder[Prepared]
	Run(...context.Context) error
}

// ---

type Runner interface {
	Run(...context.Context) error
	Raw() (string, []any, error)
}

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
