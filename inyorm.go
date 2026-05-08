package inyorm

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/impl/expression"
	"github.com/laacin/inyorm/internal/impl/statement/dml"
)

type Engine struct {
	dialect  Dialect
	instance *sql.DB
}

func New(dialect Dialect) *Engine {
	return &Engine{
		dialect:  dialect,
		instance: nil,
	}
}

func NewWithInstance(dialect Dialect, db *sql.DB) *Engine {
	return &Engine{
		dialect:  dialect,
		instance: db,
	}
}

func (eng *Engine) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := dml.NewSelectStatement(ctx, eng.dialect, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := dml.NewInsertStatement(ctx, eng.dialect, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := dml.NewUpdateStatement(ctx, eng.dialect, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := dml.NewDeleteStatement(ctx, eng.dialect, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}
