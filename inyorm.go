package inyorm

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/expression"
	"github.com/laacin/inyorm/internal/impl/statement/dml"
)

type DB struct{ eng *entity.Engine }

func New(eng *Engine) (*DB, error) {
	if eng.Err != nil {
		return nil, eng.Err
	}
	return &DB{eng}, nil
}

// --- DML Statements

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

// --- Connection

func (db *DB) Close(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() { errCh <- db.eng.Driver.Close() }()

	select {
	case <-ctx.Done():
		return errors.New("context timeout")
	case err := <-errCh:
		return err
	}
}
