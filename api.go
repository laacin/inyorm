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
	Table(name string) Table
	Col(name string, table ...string) Column
	All(table ...string) Column
	Lazy(id ...string) Parameter
	Param(value ...any) Parameter
	Cond(ident any) Condition
	Concat(v ...any) Column
	Switch(cond any, fn func(cs Case)) Column
	Search(fn func(cs Case)) Column
}
