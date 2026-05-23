package api

import "context"

type Statement interface {
	Runner
	Prepare() PrepStatement
	Bind(...any) Statement
}

type PrepStatement interface {
	BindPrep(...any) PrepStatement
	Values(...any) Runner
}

type Runner interface {
	Raw() (string, []any, error)
	Run(...context.Context) error
}
