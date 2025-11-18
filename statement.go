package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type SelectStatement struct {
	ctx   context.Context
	query *writer.Query
	*selectCls
	*fromCls
	*joinCls
	*whereCls
	*groupByCls
	*havingCls
	*orderByCls
	*limitCls
	*offsetCls
}

func (stmt *SelectStatement) Build() (string, []any) {
	clauses := []core.Clause{
		stmt.selectWrap, stmt.fromWrap, stmt.joinWrap,
		stmt.whereWrap, stmt.groupByWrap, stmt.havingWrap,
		stmt.orderByWrap, stmt.limitWrap, stmt.offsetWrap,
	}

	stmt.query.SetClauses(clauses, writer.SelectOrder)
	return stmt.query.Build()
}
