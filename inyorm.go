package inyorm

import (
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type (
	ColExpr = core.ColExpr
	Column  = core.Column
)

type QueryEngine struct {
	Dialect string
}

func (qe *QueryEngine) New(baseTable string) (*Query, ColExpr) {
	stmt := writer.NewStatement(qe.Dialect, baseTable)
	c := &column.ColExpr{Statement: stmt}
	return &Query{stmt: stmt}, c
}

type Query struct {
	stmt      *writer.Statement
	selectCls *clause.SelectClause
	joinCls   *clause.JoinClause
	whereCls  *clause.WhereClause
	groupCls  *clause.GroupByClause
	orderCls  *clause.OrderByClause
	limitCls  *clause.LimitClause
	offsetCls *clause.OffsetClause
}

func (q *Query) Distinct() core.ClauseSelect {
	cls := lazyInit(&q.selectCls)
	cls.Distinct()
	return cls
}

func (q *Query) Select(values ...Value) {
	cls := lazyInit(&q.selectCls)
	cls.Select(values...)
}

func (q *Query) Join(table string) core.ClauseOn {
	cls := lazyInit(&q.joinCls)
	return cls.Join(table)
}

func (q *Query) Where(identifier Value) core.Cond {
	cls := lazyInit(&q.whereCls)
	return cls.Where(identifier)
}

func (q *Query) GroupBy(values ...Value) core.ClauseHaving {
	cls := lazyInit(&q.groupCls)
	return cls.GroupBy(values...)
}

func (q *Query) OrderBy(value Value) core.ClauseOrder {
	cls := lazyInit(&q.orderCls)
	return cls.OrderBy(value)
}

func (q *Query) Limit(value int) *Query {
	cls := lazyInit(&q.limitCls)
	cls.Limit(value)
	return q
}

func (q *Query) Offset(value int) *Query {
	cls := lazyInit(&q.offsetCls)
	cls.Offset(value)
	return q
}

func (q *Query) Build() (string, []any) {
	from := &clause.FromClause{}
	from.From(q.stmt.GetFrom())

	clauses := []core.Clause{
		q.selectCls,
		from,
		q.joinCls,
		q.whereCls,
		q.groupCls,
		q.orderCls,
		q.limitCls,
		q.offsetCls,
	}

	q.stmt.SetClauses(clauses)

	return q.stmt.Build()
}

func lazyInit[T any](target **T) *T {
	if *target == nil {
		*target = new(T)
	}
	return *target
}
