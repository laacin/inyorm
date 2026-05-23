package query

import (
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
)

type QueryKind int

const (
	QuerySelect QueryKind = iota
	QueryInsert
	QueryUpdate
	QueryDelete

	QueryCreateTable
	QueryCreateIndex

	QueryAlterTable

	QueryDropTable
	QueryDropIndex
)

type Dialect interface {
	ddl.TableWriter
	dml.ClauseWriter
	dml.QueryWriter
	expr.ExprWriter
}

type QueryResult struct {
	Query  string
	Values []any
}

type QueryBuilder interface {
	Kind() QueryKind
	Build() (*QueryResult, error)
}
