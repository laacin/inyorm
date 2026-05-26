package query

import (
	"github.com/laacin/inyorm/internal/builder"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query/ddl"
	"github.com/laacin/inyorm/internal/query/dml"
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
	Build(*builder.Builder) error
	Render(core.InternalWriter, Dialect) error
}

type Query[T QueryApi] struct {
	dial    Dialect
	builder *builder.Builder

	API  T
	Expr *builder.ExprBuilder
}

func New[T QueryApi](api T, dial Dialect) *Query[T] {
	b := builder.New(dial)

	return &Query[T]{
		dial:    dial,
		builder: b,
		API:     api,
		Expr:    builder.NewExprBuilder(b.Params(), b.Aliases()),
	}
}

// --- Builder

func (q *Query[T]) Build() (*builder.Builder, error) {
	q.API.Build(q.builder)
	q.API.Render(q.builder.Writer(), q.dial)

	return q.builder, nil
}
