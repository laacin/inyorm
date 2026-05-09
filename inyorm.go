package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
	"github.com/laacin/inyorm/internal/impl/statement/dml"
)

type Engine struct {
	dialect Dialect
	driver  entity.Driver
}

// func New(dialect Dialect) *Engine {
// 	return &Engine{
// 		dialect: dialect,
// 		driver:  nil,
// 	}
// }

func New(driver entity.Driver, dialect Dialect) *Engine {
	return &Engine{
		dialect: dialect,
		driver:  driver,
	}
}

func (eng *Engine) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := dml.NewSelectStatement(ctx, eng.dialect, eng.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := dml.NewInsertStatement(ctx, eng.dialect, eng.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := dml.NewUpdateStatement(ctx, eng.dialect, eng.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (eng *Engine) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := dml.NewDeleteStatement(ctx, eng.dialect, eng.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}
