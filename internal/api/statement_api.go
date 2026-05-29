package api

import "context"

type Statement interface {
	Runner
	Binder[Statement]
}

type Binder[T any] interface {
	Bind(any) T
	Values(any) T
}

type SelfBinder interface {
	Bind(any) SelfBinder
	Values(any) SelfBinder
}

type Runner interface {
	Raw() (string, []any, error)
	Run(...context.Context) error
}
