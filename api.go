package inyorm

import "context"

// --- Main input

type InyormBuilder interface {
	// Dml
	Select(func(q SelectQuery, e Expr)) Statement
	Insert(func(q InsertQuery, e Expr)) Statement
	Update(func(q UpdateQuery, e Expr)) Statement
	Delete(func(q DeleteQuery, e Expr)) Statement

	// Execution
	RunTx(context.Context, ...Runner) error
}

// --- Queries

type Expr interface {
	Table(name string) any
	Col(name string, table ...string) Column
	All(table ...string) Column
	// Lazy(ref ...any) Parameter
	Param(value ...any) any
	Cond(ident any) Condition
	Concat(v ...any) Column
	Switch(cond any, fn func(cs Case)) Column
	Search(fn func(cs Case)) Column
}
