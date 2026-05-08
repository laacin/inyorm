package clause

import (
	"slices"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/mapper"
)

type InsertIntoImpl struct {
	declared bool
	emb      entity.InsertInto
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

func (c *InsertIntoImpl) Kind() entity.ClauseKind {
	return entity.ClauseInsertInto
}

func (c *InsertIntoImpl) Build() entity.Clause {
	if len(c.ref) < 1 {
		panic("TODO")
		// return errors.New("missing reference")
	}

	cols, err := mapper.GetColumns("inyorm", c.ref)
	if err != nil {
		panic("TODO")
		// return fmt.Errorf("failed to get columns: %w", err)
	}

	ignores, err := mapper.GetColumns("inyorm", c.ignores)
	if err != nil {
		panic("TODO")
		// return fmt.Errorf("failed to get columns: %w", err)
	}

	cols = slices.DeleteFunc(cols, func(col string) bool {
		return slices.Contains(ignores, col)
	})

	result, err := mapper.Read("inyorm", cols, c.values)
	if err != nil {
		panic("TODO")
		// return fmt.Errorf("failed to map value: %w", err)
	}

	params := make([]any, len(result.Args))
	for i, arg := range result.Args {
		params[i] = &entity.Parameter{
			Store: true,
			Value: arg,
		}
	}

	c.emb.Cols = cols
	c.emb.Rows = result.Rows
	c.emb.Values = params
	return &c.emb
}
