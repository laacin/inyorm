package inyorm

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal/core"
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

// --- Transaction

func (db *DB) RunTx(ctx context.Context, stmts ...Runner) error {
	tx := db.eng.Driver.BeginTx(ctx)

	for _, stmt := range stmts {
		stmt, ok := stmt.(interface {
			RunTx(context.Context, core.Transaction) error
		})
		if !ok {
			_ = tx.Rollback()
			return errors.New("runner does not support transactions")
		}

		if err := stmt.RunTx(ctx, tx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// --- DML Statements

func (db *DB) Select(ref string, fn func(q SelectQuery, e Expr)) Statement {
	q := &statement.SelectQueryImpl{}
	e := &exprimpl.ExprBuilderImpl{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &statement.StatementImpl{}
	return stmt.Start(db.eng.Driver, q)
}
func (db *DB) Insert(ref string, fn func(q InsertQuery, e Expr)) Statement {
	q := &statement.InsertQueryImpl{}
	e := &exprimpl.ExprBuilderImpl{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &statement.StatementImpl{}
	return stmt.Start(db.eng.Driver, q)
}
func (db *DB) Update(ref string, fn func(q UpdateQuery, e Expr)) Statement {
	q := &statement.UpdateQueryImpl{}
	e := &exprimpl.ExprBuilderImpl{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &statement.StatementImpl{}
	return stmt.Start(db.eng.Driver, q)
}
func (db *DB) Delete(ref string, fn func(q DeleteQuery, e Expr)) Statement {
	q := &statement.DeleteQueryImpl{}
	e := &exprimpl.ExprBuilderImpl{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &statement.StatementImpl{}
	return stmt.Start(db.eng.Driver, q)
}

// --- DDL Statements

func (db *DB) CreateTable(name string, fn func(q CreateTable, e Expr)) Statement {
	q := &statement.CreateTableQueryImpl{}
	e := &exprimpl.ExprBuilderImpl{}

	fn(q.Start(db.eng.Dialect, name), e.Start(name))
	stmt := &statement.StatementImpl{}
	return stmt.Start(db.eng.Driver, q)
}

// --- Connection

func (db *DB) Close() error {
	return db.eng.Driver.Close()
}
