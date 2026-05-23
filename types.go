package inyorm

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
)

// --- Main types
type Engine struct {
	Dialect query.Dialect
	Driver  core.Driver
}

// --- Expressions
type (
	// Table         = api.Table
	Column = api.Col
	// Parameter     = api.Parameter
	Condition     = api.Cond
	ConditionNext = api.CondNext
	Case          = api.Case
	CaseNext      = api.CaseNext
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
