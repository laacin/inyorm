package inyorm

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir"
)

// --- Main types
type Engine = ir.Engine

// --- Expressions
type (
	Parameter     = api.Parameter
	Table         = api.Table
	Column        = api.Column
	Condition     = api.Condition
	ConditionNext = api.ConditionNext
	Case          = api.Case
	CaseNext      = api.CaseNext

	// Builder
	ExprBuilder = api.ExprBuilder[Case]
)

// --- Executor
type (
	Prepare  = api.Prepare
	Executor = api.Executor
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

	// Statements
	SelectStatement = api.SelectStmt
	InsertStatement = api.InsertStmt
	UpdateStatement = api.UpdateStmt
	DeleteStatement = api.DeleteStmt
)

// --- DDL
type (
	CreateTable = api.CreateTable
	CreateIndex = api.CreateIndex
)
