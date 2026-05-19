package inyorm

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/impl/statement"
	"github.com/laacin/inyorm/internal/impl/table"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
)

type DB struct{ eng *ir.Engine }

func New(eng *Engine) (*DB, error) {
	if eng.Err != nil {
		return nil, eng.Err
	}
	return &DB{eng}, nil
}

// --- DDL Statements

func (db *DB) NewTable(ctx context.Context, name string, fn func(tb TableBuilder, e ExprBuilder)) string {
	tb := &table.TableBuilderImpl{}
	fn(tb.Start(name), &exprimpl.ExprBuilderImpl{DefaultRef: name})

	w := &writer.WriterImpl{Syntax: db.eng.Dialect}

	tb.Build(w, db.eng.Dialect)
	return w.ToString()
}

// --- DML Statements

func (db *DB) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := statement.NewSelectStatement(ctx, db.eng, table)
	exprBuilder := &exprimpl.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := statement.NewInsertStatement(ctx, db.eng, table)
	exprBuilder := &exprimpl.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := statement.NewUpdateStatement(ctx, db.eng, table)
	exprBuilder := &exprimpl.ExprBuilderImpl{DefaultRef: table}
	return stmt, exprBuilder
}

func (db *DB) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := statement.NewDeleteStatement(ctx, db.eng, table)
	exprBuilder := &exprimpl.ExprBuilderImpl{DefaultRef: table}
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
