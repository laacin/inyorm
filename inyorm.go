package inyorm

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/writer"
)

type QueryEngine struct {
	dialect string
}

func New(dialect string, db *sql.DB) *QueryEngine {
	qe := &QueryEngine{dialect: dialect}

	return qe
}

// ---- Statements

func (qe *QueryEngine) NewSelect(ctx context.Context, table string) (*SelectStatement, *ColumnExpr) {
	query := writer.NewQuery(qe.dialect, table)
	stmt := &SelectStatement{
		query:      query,
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
	stmt.From(table)

	return stmt, newColExpr(table)
}
