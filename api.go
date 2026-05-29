package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
)

// --- Queries

type Queries[T any] interface {
	Select(func(q SelectQuery, e Expr)) T
	Insert(func(q InsertQuery, e Expr)) T
	Update(func(q UpdateQuery, e Expr)) T
	Delete(func(q DeleteQuery, e Expr)) T

	CreateTable(func(q CreateTable, e Expr)) T
}

// --- ExprBuilder

type Expr interface {
	Table(name string) Table
	Col(name string, table ...string) Col
	All(table ...string) Col
	Param(value any) Param
	Cond(ident any) Cond
	Concat(v ...any) Col
	Switch(cond any, fn func(cs Case)) Col
	Search(fn func(cs Case)) Col
}

// ---

type Transaction interface {
	Queries[api.SelfBinder]
	Run(...context.Context) error
}
