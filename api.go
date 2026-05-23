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
