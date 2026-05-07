package inyorm

import "context"

type Engine struct{ dialect Dialect }

func New(dialect Dialect) *Engine {
	return &Engine{dialect}
}

func (eng *Engine) NewSelect(ctx context.Context, table string) (SelectStmt, ExprBuilder) {
	stmt := newSelect(ctx, eng.dialect, table)
	exprBuilder := &exprBuilder{Dialect: eng.dialect, MainRef: table}
	return stmt, exprBuilder
}
