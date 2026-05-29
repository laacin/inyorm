package inyorm

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
)

// --- Main types
type Driver = core.Driver

type Engine struct {
	Dialect Dialect
	Driver  Driver
}

type Dialect interface {
	expr.Renderer
	dml.Renderer
	ddl.Renderer
}

// --- Expressions
type (
	Table    = any
	Col      = api.Col
	Param    = any
	Cond     = api.Cond
	CondNext = api.CondNext
	Case     = api.Case
	CaseNext = api.CaseNext
)

// --- Statement
type (
	Statement = api.Statement
	Runner    = api.Runner
)

// --- DML
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

// --- DDL
type (
	CreateTable = api.CreateTableQuery
)
