package inyorm

import (
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

func (db *DB) CreateTable(name string, fn func(q CreateTable, e Expr)) Statement {
	q := &statement.CreateTableQueryImpl{}
	e := &exprimpl.ExprBuilderImpl{}

	fn(q.Start(db.eng.Dialect, name), e.Start(name))
	stmt := &statement.StatementImpl{}
	return stmt.Start(db.eng.Driver, q)
}

func (db *DB) Close() error {
	return db.eng.Driver.Close()
}
