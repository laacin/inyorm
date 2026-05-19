package inyorm

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/impl/statement"
	"github.com/laacin/inyorm/internal/ir"
)

type DB struct{ eng *ir.Engine }

func New(eng *Engine) (*DB, error) {
	if eng.Err != nil {
		return nil, eng.Err
	}
	return &DB{eng}, nil
}

// --- DML Statements

func (db *DB) NewSelect(ctx context.Context, table string) (SelectStatement, ExprBuilder) {
	stmt := &statement.SelectStmtImpl{}
	e := &exprimpl.ExprBuilderImpl{}
	return stmt.Start(ctx, db.eng, table), e.Start(table)
}

func (db *DB) NewInsert(ctx context.Context, table string) (InsertStatement, ExprBuilder) {
	stmt := &statement.InsertStmtImpl{}
	e := &exprimpl.ExprBuilderImpl{}
	return stmt.Start(ctx, db.eng, table), e.Start(table)
}

func (db *DB) NewUpdate(ctx context.Context, table string) (UpdateStatement, ExprBuilder) {
	stmt := &statement.UpdateStmtImpl{}
	e := &exprimpl.ExprBuilderImpl{}
	return stmt.Start(ctx, db.eng, table), e.Start(table)
}

func (db *DB) NewDelete(ctx context.Context, table string) (DeleteStatement, ExprBuilder) {
	stmt := &statement.DeleteStmtImpl{}
	e := &exprimpl.ExprBuilderImpl{}
	return stmt.Start(ctx, db.eng, table), e.Start(table)
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
