package clause

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/impl/mapper"
	"github.com/laacin/inyorm/internal/ir/dml"
)

type UpdateImpl struct {
	declared bool
	emb      dml.Update
	ref      []any
	ignores  []any
	values   any
}

func (c *UpdateImpl) Update(reference ...any) api.Values {
	c.declared = true
	c.ref = reference
	return c
}

func (c *UpdateImpl) UpdateIgnore(reference any, ignores ...any) api.Values {
	c.declared = true
	c.ref = []any{reference}
	c.ignores = ignores
	return c
}

func (c *UpdateImpl) Values(values any) {
	c.values = values
}

func (c *UpdateImpl) Table(table any) {
	c.emb.Table = table
}

// --- Build

func (c *UpdateImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *UpdateImpl) Kind() dml.ClauseKind {
	return dml.ClauseUpdate
}

func (c *UpdateImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
	if len(c.ref) < 1 {
		return errors.New("missing reference")
	}

	cols := mapper.ReadColumns(c.ref)
	ignores := mapper.ReadColumns(c.ignores)

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	result, err := mapper.ReadValues(cols, c.values)
	if err != nil {
		return fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(result.Args))
	for i, arg := range result.Args {
		params[i] = (&exprimpl.ParameterImpl{}).Start(true, arg)
	}

	c.emb.Cols = cols
	c.emb.Values = params

	dial.WriteUpdate(w, &c.emb)
	return nil
}
