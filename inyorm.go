package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
	"github.com/laacin/inyorm/internal/impl/statement/dml"
)

type DB struct {
	dialect Dialect
	driver  entity.Driver
}

func New(driver entity.Driver, dialect Dialect) *DB {
	return &DB{
		dialect: dialect,
		driver:  driver,
	}
}

func (db *DB) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := dml.NewSelectStatement(ctx, db.dialect, db.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := dml.NewInsertStatement(ctx, db.dialect, db.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := dml.NewUpdateStatement(ctx, db.dialect, db.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := dml.NewDeleteStatement(ctx, db.dialect, db.driver, table)
	exprBuilder := &expression.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}
