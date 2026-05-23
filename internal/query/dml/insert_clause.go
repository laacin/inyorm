package dml

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/impl/mapper"
)

// --- Entity

type Insert struct {
	Table  any
	Cols   []string
	Rows   int
	Values []any
}

// --- Builder

type InsertBuilder struct {
	declared bool
	emb      Insert
	ref      []any
	ignores  []any
	values   any
}

// --- PUB API

func (b *InsertBuilder) Insert(ref ...any) api.Ignore {
	b.declared = true
	b.ref = ref
	return b
}

func (b *InsertBuilder) Ignore(ignore ...any) {
	b.ignores = ignore
}

func (b *InsertBuilder) Values(v any) {
	b.values = v
}

func (b *InsertBuilder) Into(v any) {
	b.emb.Table = v
}

// --- Build

func (*InsertBuilder) Kind() ClauseKind {
	return ClauseInsert
}

func (b *InsertBuilder) IsDeclared() bool {
	return b != nil && b.declared
}

func (b *InsertBuilder) Build(w core.InternalWriter, dial ClauseWriter) error {
	if len(b.ref) < 1 {
		return errors.New("missing reference")
	}

	cols := mapper.ReadColumns(b.ref)
	ignores := mapper.ReadColumns(b.ignores)

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	result, err := mapper.ReadValues(cols, b.values)
	if err != nil {
		return fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(result.Args))
	for i, arg := range result.Args {
		params[i] = (&expr.ParamBuilder{}).Start(true, arg)
	}

	b.emb.Cols = cols
	b.emb.Rows = result.Rows
	b.emb.Values = params

	dial.WriteInsertInto(w, &b.emb)
	return nil
}
