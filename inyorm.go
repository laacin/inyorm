package inyorm

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
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

func (db *DB) Select(fn func(q SelectQuery, e Expr)) Statement {
	q := dml.NewSelect(db.eng.Dialect)

	qc := query.NewCompiler(q)
	fn(q, qc.Expr())

	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, db.eng.Dialect, qc)
}
func (db *DB) Insert(fn func(q InsertQuery, e Expr)) Statement {
	q := dml.NewInsert(db.eng.Dialect)

	qc := query.NewCompiler(q)
	fn(q, qc.Expr())

	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, db.eng.Dialect, qc)
}
func (db *DB) Update(fn func(q UpdateQuery, e Expr)) Statement {
	q := dml.NewUpdate(db.eng.Dialect)

	qc := query.NewCompiler(q)
	fn(q, qc.Expr())

	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, db.eng.Dialect, qc)
}
func (db *DB) Delete(fn func(q DeleteQuery, e Expr)) Statement {
	q := dml.NewDelete(db.eng.Dialect)

	qc := query.NewCompiler(q)
	fn(q, qc.Expr())

	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, db.eng.Dialect, qc)
}

// --- DDL Statements

func (db *DB) CreateTable(fn func(q CreateTable, e Expr)) Statement {
	q := ddl.NewCreateTable(db.eng.Dialect)

	qc := query.NewCompiler(q)
	fn(q, qc.Expr())

	stmt := &internal.Statement{}
	return stmt.Start(db.eng.Driver, db.eng.Dialect, qc)
}

// --- Connection

func (db *DB) Close() error {
	return db.eng.Driver.Close()
}
