package inyorm

import (
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
)

type queryBuilder[T any] struct {
	dial  Dialect
	chain func(qc *query.Compiler) T
}

func newQueryBuilder[T any](dial Dialect, chain func(qc *query.Compiler) T) Queries[T] {
	return &queryBuilder[T]{
		dial:  dial,
		chain: chain,
	}
}

func (qb *queryBuilder[T]) CreateTable(fn func(q CreateTable, e Expr)) T {
	q := ddl.NewCreateTable(qb.dial)

	qc := query.NewCompiler(q, expr.NewParser(qb.dial))
	fn(q, qc.ExprBuilder())

	return qb.chain(qc)
}

// --- DML

func (qb *queryBuilder[T]) Select(fn func(q SelectQuery, e Expr)) T {
	q := dml.NewSelect(qb.dial)

	qc := query.NewCompiler(q, expr.NewParser(qb.dial))
	fn(q, qc.ExprBuilder())

	return qb.chain(qc)
}
func (qb *queryBuilder[T]) Insert(fn func(q InsertQuery, e Expr)) T {
	q := dml.NewInsert(qb.dial)

	qc := query.NewCompiler(q, expr.NewParser(qb.dial))
	fn(q, qc.ExprBuilder())

	return qb.chain(qc)
}
func (qb *queryBuilder[T]) Update(fn func(q UpdateQuery, e Expr)) T {
	q := dml.NewUpdate(qb.dial)

	qc := query.NewCompiler(q, expr.NewParser(qb.dial))
	fn(q, qc.ExprBuilder())

	return qb.chain(qc)
}
func (qb *queryBuilder[T]) Delete(fn func(q DeleteQuery, e Expr)) T {
	q := dml.NewDelete(qb.dial)

	qc := query.NewCompiler(q, expr.NewParser(qb.dial))
	fn(q, qc.ExprBuilder())

	return qb.chain(qc)
}
