package api

import "context"

type Statement interface {
	Runner
	Binder[Statement]
	Prepare(...context.Context) (Prepared, error)
}

type Prepared interface {
	Run(...context.Context) error
	Bind(any) Prepared
	Values(any, ...string) Prepared
}

type Binder[T any] interface {
	Bind(any) T
	Values(any, ...string) T
}

type SelfBinder interface {
	Bind(any) SelfBinder
	Values(any, ...string) SelfBinder
}

type Runner interface {
	Raw() (string, []any, error)
	Run(...context.Context) error
}
