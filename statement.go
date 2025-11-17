package inyorm

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type SelectStatement struct {
	builder *writer.StatementBuilder
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
	stmt.builder.SetClauses(clauses, writer.SelectOrder)
	return stmt.builder.Build()
}
