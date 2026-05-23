package dml

import (
	"errors"
	"fmt"
	"slices"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/expr"
	"github.com/laacin/inyorm/internal/impl/mapper"
)

// --- Entity

type ClauseInsert struct {
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

func (c *ClauseInsert) Insert(ref ...any) api.Ignore {
	c.declared = true
	c.ref = ref
	return c
}

func (c *ClauseInsert) Ignore(ignore ...any) {
	c.ignores = ignore
}

func (c *ClauseInsert) Values(v any) {
	c.rawVal = v
}

func (c *ClauseInsert) Into(v any) {
	c.Table = v
}

// --- Build

func (*ClauseInsert) Kind() ClauseKind {
	return ClauseKindInsert
}

func (c *ClauseInsert) IsDeclared() bool {
	return c != nil && c.declared
}

func (c *ClauseInsert) Build() error {
	if len(c.ref) < 1 {
		return errors.New("missing reference")
	}

	cols := mapper.ReadColumns(c.ref)
	ignores := mapper.ReadColumns(c.ignores)

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	result, err := mapper.ReadValues(cols, c.rawVal)
	if err != nil {
		return fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(result.Args))
	for i, arg := range result.Args {
		params[i] = (&expr.Param{}).Start(true, arg)
	}

	c.Cols = cols
	c.Rows = result.Rows
	c.Vals = params

	c.ref = nil
	c.ignores = nil
	c.rawVal = nil

	return nil
}
