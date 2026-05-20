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

type InsertIntoImpl struct {
	declared bool
	emb      dml.InsertInto
	ref      []any
	ignores  []any
	values   any
}

func (c *InsertIntoImpl) Insert(reference ...any) api.Ignore {
	c.declared = true
	c.ref = reference
	return c
}

func (c *InsertIntoImpl) Ignore(ignores ...any) {
	c.ignores = ignores
}

func (c *InsertIntoImpl) Values(values any) {
	c.values = values
}

func (c *InsertIntoImpl) Into(table any) {
	c.emb.Table = table
}

// --- Build

func (c *InsertIntoImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *InsertIntoImpl) Kind() dml.ClauseKind {
	return dml.ClauseInsertInto
}

func (c *InsertIntoImpl) Build(w core.InternalWriter, dial dml.ClauseWriter) error {
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
	c.emb.Rows = result.Rows
	c.emb.Values = params

	dial.WriteInsertInto(w, &c.emb)
	return nil
}
