package clause

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/mapper"
	"github.com/laacin/inyorm/internal/ir/dml"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type InsertIntoImpl struct {
	declared bool
	emb      dml.InsertInto
	ref      []any
	ignores  []any
	values   any
}

func (c *InsertIntoImpl) Insert(reference ...any) api.Values {
	c.declared = true
	c.ref = reference
	return c
}

func (c *InsertIntoImpl) InsertIgnore(reference any, ignores ...any) api.Values {
	c.declared = true
	c.ref = []any{reference}
	c.ignores = ignores
	return c
}

func (c *InsertIntoImpl) Values(values any) {
	c.values = values
}

func (c *InsertIntoImpl) Table(table any) {
	c.emb.Table = table
}

// --- Build

func (c *InsertIntoImpl) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *InsertIntoImpl) Kind() dml.ClauseKind {
	return dml.ClauseInsertInto
}

func (c *InsertIntoImpl) Build() (dml.Clause, error) {
	if len(c.ref) < 1 {
		return nil, errors.New("missing reference")
	}

	cols := mapper.ReadColumns(c.ref)
	ignores := mapper.ReadColumns(c.ignores)

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	result, err := mapper.ReadValues(cols, c.values)
	if err != nil {
		return nil, fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(result.Args))
	for i, arg := range result.Args {
		params[i] = &expr.Parameter{
			Store: true,
			Value: arg,
		}
	}

	c.emb.Cols = cols
	c.emb.Rows = result.Rows
	c.emb.Values = params
	return &c.emb, nil
}
