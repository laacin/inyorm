package inyorm

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type (
	Engine = entity.Engine

	// Values
	Value   = api.Value
	Scanner = api.Scanner

	Parameter     = api.Parameter
	Table         = api.Table
	Column        = api.Column
	Condition     = api.Condition
	ConditionNext = api.ConditionNext
	Case          = api.Case
	CaseNext      = api.CaseNext

	// Executor
	Prepare  = api.Prepare
	Executor = api.Executor

	// Expression Builder
	ExprBuilder = api.ExprBuilder[Case]

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
