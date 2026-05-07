package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/impl/expression"
)

type Engine struct{ dialect Dialect }

func New(dialect Dialect) *Engine {
	return &Engine{dialect}
}

func (eng *Engine) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := newSelect(ctx, eng.dialect, table)
	exprBuilder := &expression.ExpressionImpl{Dialect: eng.dialect, MainRef: table}
	return stmt, exprBuilder
}
