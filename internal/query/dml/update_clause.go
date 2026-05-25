package dml

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

// --- Entity

type ClauseUpdate struct {
	Table any
	Cols  []string
	Vals  []any

	// internal
	declared bool
	ref      []any
	ignores  []any
	rawVal   any
}

// --- PUB API

func (c *ClauseUpdate) Update(ref ...any) api.Ignore {
	c.declared = true
	c.ref = ref
	return c
}

func (c *ClauseUpdate) Ignore(ignore ...any) {
	c.ignores = ignore
}

func (c *ClauseUpdate) Values(v any) {
	c.rawVal = v
}

func (c *ClauseUpdate) Into(v any) {
	c.Table = v
}

// --- Build

func (*ClauseUpdate) Kind() ClauseKind {
	return ClauseKindUpdate
}

func (c *ClauseUpdate) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseUpdate) Build(b *core.Builder) error {
	if len(c.ref) < 1 {
		return errors.New("missing reference")
	}

	cols := b.Mapper.ReadCols(c.ref...)
	ignores := b.Mapper.ReadCols(c.ignores...)

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	if ph, ok := c.rawVal.(*expr.Placeholder); ok && ph.IsLazy() {
		params := make([]any, len(cols))
		for i := range cols {
			ph := &expr.Placeholder{}
			if i == 0 {
				ph.Start(func() core.ParamIndex {
					b.Param.LazyObj(cols)
					return b.Param.LastIndex(len(cols) - 1)
				})
			}
			params[i] = ph.Start(func() core.ParamIndex {
				return b.Param.LastIndex(len(cols) - i - 1)
			})
		}

		return c.done(cols, params)
	}

	args, err := b.Mapper.ReadValues(cols, c.rawVal)
	if err != nil {
		return fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(args))
	for i, arg := range args {
		params[i] = (&expr.Placeholder{}).Start(func() core.ParamIndex {
			b.Param.Store(arg)
			return b.Param.LastIndex(0)
		})
	}

	return c.done(cols, params)
}

// helpers
func (c *ClauseUpdate) done(cols []string, vals []any) error {
	c.Cols = cols
	c.Vals = vals

	c.ref = nil
	c.ignores = nil
	c.rawVal = nil

	return nil
}
