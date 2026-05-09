package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
	"github.com/laacin/inyorm/internal/impl/statement/dml"
)

type Engine struct {
	Dialect entity.Dialect
	Driver  entity.Driver
}

type DB struct {
	eng *Engine
}

func New(eng *Engine) *DB {
	return &DB{eng}
}

func (db *DB) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := dml.NewSelectStatement(ctx, db.eng.Dialect, db.eng.Driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := dml.NewInsertStatement(ctx, db.eng.Dialect, db.eng.Driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := dml.NewUpdateStatement(ctx, db.eng.Dialect, db.eng.Driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := dml.NewDeleteStatement(ctx, db.eng.Dialect, db.eng.Driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}
