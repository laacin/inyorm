package dml

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/mapper"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/query"
)

// --- Entity

type ClauseInsertInto struct {
	Table any
	Cols  []string
	Rows  int
	Vals  []any

	// internal
	declared bool
	ref      []any
	ignores  []any
	rawVal   any
}

// --- PUB API

func (c *ClauseInsertInto) Insert(ref ...any) api.Ignore {
	c.declared = true
	c.ref = ref
	return c
}

func (c *ClauseInsertInto) Ignore(ignore ...any) {
	c.ignores = ignore
}

func (c *ClauseInsertInto) Values(v any) {
	c.rawVal = v
}

func (c *ClauseInsertInto) Into(v any) {
	c.Table = v
}

// --- Build

func (*ClauseInsertInto) Kind() ClauseKind {
	return ClauseKindInsertInto
}

func (c *ClauseInsertInto) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseInsertInto) Build(tools *query.Tools) error {
	if len(c.ref) < 1 {
		return errors.New("missing reference")
	}

	m := mapper.New()

	cols := m.ReadCols(c.ref...)
	ignores := m.ReadCols(c.ignores...)

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	if ph, ok := c.rawVal.(*expr.Placeholder); ok && ph.IsLazy() {
		params := make([]any, len(cols))

		for i := range cols {
			ph := &expr.Placeholder{}

			if i == 0 {
				params[i] = ph.Start(func() core.ParamIndex {
					tools.Params.LazyObject(cols)
					return tools.Params.LastIndex(len(cols) - 1)
				})
				continue
			}

			params[i] = ph.Start(func() core.ParamIndex {
				return tools.Params.LastIndex(len(cols) - i - 1)
			})
		}

		return c.done(cols, 1, params)
	}

	args, err := m.ReadValues(cols, c.rawVal)
	if err != nil {
		return fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(args))
	for i, arg := range args {
		params[i] = (&expr.Placeholder{}).Start(func() core.ParamIndex {
			tools.Params.Store(arg)
			return tools.Params.LastIndex(0)
		})
	}

	return c.done(cols, len(args)/len(cols), params)
}

// helpers
func (c *ClauseInsertInto) done(cols []string, rows int, vals []any) error {
	c.Cols = cols
	c.Rows = rows
	c.Vals = vals

	c.ref = nil
	c.ignores = nil
	c.rawVal = nil

	return nil
}
