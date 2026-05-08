package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/impl/expression"
	"github.com/laacin/inyorm/internal/impl/statement/dml"
)

type Engine struct{ dialect Dialect }

func New(dialect Dialect) *Engine {
	return &Engine{dialect}
}

func (eng *Engine) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := &dml.SelectStmtImpl{Dialect: eng.dialect, DefaultRef: table}
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := &dml.InsertStmtImpl{Dialect: eng.dialect, DefaultRef: table}
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := &dml.UpdateStmtImpl{Dialect: eng.dialect, DefaultRef: table}
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := &dml.DeleteStmtImpl{Dialect: eng.dialect, DefaultRef: table}
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}
