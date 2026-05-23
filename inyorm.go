package inyorm

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
)

type DB struct{ eng *Engine }

func New(eng *Engine, err error) (*DB, error) {
	if err != nil {
		return nil, err
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
	q := &query.SelectQuery{}
	e := &internal.ExprBuilder{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, q)
}
func (db *DB) Insert(ref string, fn func(q InsertQuery, e Expr)) Statement {
	q := &query.InsertQuery{}
	e := &internal.ExprBuilder{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, q)
}
func (db *DB) Update(ref string, fn func(q UpdateQuery, e Expr)) Statement {
	q := &query.UpdateQuery{}
	e := &internal.ExprBuilder{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, q)
}
func (db *DB) Delete(ref string, fn func(q DeleteQuery, e Expr)) Statement {
	q := &query.DeleteQuery{}
	e := &internal.ExprBuilder{}

	fn(q.Start(db.eng.Dialect, ref), e.Start(ref))
	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, q)
}

// --- DDL Statements

func (db *DB) CreateTable(name string, fn func(q CreateTable, e Expr)) Statement {
	q := &query.CreateTableQuery{}
	e := &internal.ExprBuilder{}

	fn(q.Start(db.eng.Dialect, name), e.Start(name))
	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, q)
}

// --- Connection

func (db *DB) Close() error {
	return db.eng.Driver.Close()
}
