package inyorm

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type SelectStatement struct {
	stmt *writer.Statement
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
	stmt.stmt.SetClauses(clauses)
	return stmt.stmt.Build(writer.SelectOrder)
}

func NewSelect(table string) (*SelectStatement, *ColumnExpr) {
	stmt := writer.NewStatement("", table)
	statement := &SelectStatement{
		stmt:       stmt,
		selectCls:  wrapSelect(),
		fromCls:    wrapFrom(),
		joinCls:    wrapJoin(),
		whereCls:   wrapWhere(),
		groupByCls: wrapGroupBy(),
		havingCls:  wrapHaving(),
		orderByCls: wrapOrderBy(),
		limitCls:   wrapLimit(),
		offsetCls:  wrapOffset(),
	}
	statement.From(table)

	return statement, newColExpr(table)
}
