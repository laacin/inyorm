package inyorm

import (
	"database/sql"

	"github.com/laacin/inyorm/internal/writer"
)

type QueryEngine struct {
	dialect  string
	instance *sql.DB
}

func New(dialect string, instance *sql.DB) *QueryEngine {
	return &QueryEngine{
		dialect:  dialect,
		instance: instance,
	}
}

// ---- Statements

func (qe *QueryEngine) NewSelect(table string) (*SelectStatement, *ColumnExpr) {
	stmt := writer.NewStatement(qe.dialect, table)
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
