package query

import (
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
	"github.com/laacin/inyorm/internal/writer"
)

// ---- Dialect

type Dialect interface {
	ddl.TableWriter
	dml.ClauseWriter
	expr.ExprWriter

	// DDL
	WriteQueryCreateTable(core.Writer, *QueryCreateTable)

	// DML
	WriteQuerySelect(core.Writer, *QuerySelect)
	WriteQueryInsert(core.Writer, *QueryInsert)
	WriteQueryUpdate(core.Writer, *QueryUpdate)
	WriteQueryDelete(core.Writer, *QueryDelete)
}

// --- Types

type QueryApi interface {
	Build(*core.Builder) error
	Render(core.InternalWriter, Dialect) error
}

type QueryBuilder interface {
	Build() (*QueryResult, error)
}

type Query[T QueryApi] struct {
	dial    Dialect
	builder *core.Builder
	ref     string // TODO: make it implicit

	Expr *builder.ExprBuilder
	API  T
}

func New[T QueryApi](api T, dial Dialect, ref string) *Query[T] {
	b := builder.New()
	e := (&builder.ExprBuilder{}).Start(b.Param)
	e.AttachRef(ref)

	return &Query[T]{
		ref:     ref,
		dial:    dial,
		builder: b,
		API:     api,
		Expr:    e,
	}
}

// --- Builder

func (q *Query[T]) Build() (*QueryResult, error) {
	q.API.Build(q.builder)

	w := writer.New(q.dial, q.builder.Attach.UseAliases)
	q.API.Render(w, q.dial)

	vals, err := q.builder.Param.Values()
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: vals,
	}, nil
}

// --- Result

type QueryResult struct {
	Query  string
	Values []any
}
