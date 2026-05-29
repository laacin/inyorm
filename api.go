package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
)

// --- Main types

type Engine struct {
	Dialect Dialect
	Driver  Driver
}

type Driver = core.Driver
type Dialect interface {
	expr.Renderer
	dml.Renderer
	ddl.Renderer
}

// --- Public

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

type Transaction interface {
	queries[api.SelfBinder]
	Run(...context.Context) error
}

// --- Re-exports

// Expressions
type (
	Table    = any
	Col      = api.Col
	Param    = any
	Cond     = api.Cond
	CondNext = api.CondNext
	Case     = api.Case
	CaseNext = api.CaseNext
)

// Statement
type (
	Statement = api.Statement
	Runner    = api.Runner
)

// DML
type (
	// Clauses
	Select      = api.Select
	SelectNext  = api.SelectNext
	From        = api.From
	Join        = api.Join
	JoinNext    = api.JoinNext
	JoinEnd     = api.JoinEnd
	Where       = api.Where
	GroupBy     = api.GroupBy
	Having      = api.Having
	OrderBy     = api.OrderBy
	OrderByNext = api.OrderByNext
	Limit       = api.Limit
	Offset      = api.Offset
	InsertInto  = api.Insert
	Update      = api.Update
	Delete      = api.Delete

	// Queries
	SelectQuery = api.SelectQuery
	InsertQuery = api.InsertQuery
	UpdateQuery = api.UpdateQuery
	DeleteQuery = api.DeleteQuery
)

// DDL
type (
	CreateTable = api.CreateTableQuery
)

// --- Internal

type queries[T any] interface {
	Select(func(q SelectQuery, e Expr)) T
	Insert(func(q InsertQuery, e Expr)) T
	Update(func(q UpdateQuery, e Expr)) T
	Delete(func(q DeleteQuery, e Expr)) T

	CreateTable(func(q CreateTable, e Expr)) T
}
