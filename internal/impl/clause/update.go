package clause

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/expr"
	"github.com/laacin/inyorm/internal/mapper"
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

func (c *UpdateImpl) Build() (dml.Clause, error) {
	if len(c.ref) < 1 {
		return nil, errors.New("missing reference")
	}

	cols, err := mapper.GetColumns("inyorm", c.ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	ignores, err := mapper.GetColumns("inyorm", c.ignores)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	result, err := mapper.Read("inyorm", cols, c.values)
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
	c.emb.Values = params
	return &c.emb, nil
}
