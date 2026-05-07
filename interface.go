package inyorm

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/impl/expression"
)

type (
	Dialect = entity.Dialect

	exprBuilder = expression.ExpressionImpl[
		Column, Param,
		Condition, ConditionNext,
		Case, CaseNext,
	]

	clsSelect  = clause.SelectImpl[SelectNext]
	clsFrom    = clause.FromImpl
	clsJoin    = clause.JoinImpl[JoinNext, JoinEnd, Condition, ConditionNext]
	clsWhere   = clause.WhereImpl[Condition, ConditionNext]
	clsGroupBy = clause.GroupByImpl
	clsHaving  = clause.HavingImpl[Condition, ConditionNext]
	clsOrderBy = clause.OrderByImpl[OrderByNext]
	clsLimit   = clause.LimitImpl
	clsOffset  = clause.OffsetImpl

	// TODO: insert, update, delete statements
)
