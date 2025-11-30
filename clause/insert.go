package clause

import (
	"errors"
	"fmt"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/mapper"
)

type InsertInto[Next any] struct {
	declared  bool
	table     string
	reference []any
	values    any
}

func (cls *InsertInto[Next]) Name() string     { return "INSERT INTO" }
func (cls *InsertInto[Next]) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *InsertInto[Next]) Build(w core.Writer, cfg *core.Config) error {
	if len(cls.reference) < 1 {
		return errors.New("missing reference")
	}

	ref := cls.reference[0]
	if len(cls.reference) > 1 {
		ref = cls.reference
	}

	cols, err := mapper.GetColumns(cfg.ColumnTag, ref)
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	var result *mapper.ReadResult
	if cls.values != nil {
		result, err = mapper.Read(cfg.ColumnTag, cols, cls.values)
		if err != nil {
			return fmt.Errorf("failed to map value: %w", err)
		}
	}

	var (
		ph   = false
		rows = 1
	)

	if result != nil {
		if result.Rows == 0 {
			return errors.New("there should be at least one row")
		}
		ph = true
		rows = result.Rows
	}

	w.Write("INSERT INTO")
	w.Char(' ')
	w.Table(cls.table)
	w.Char(' ')
	w.Char('(')
	for i, col := range cols {
		if i > 0 {
			w.Write(", ")
		}
		w.Column(cls.table, col)
	}

	w.Char(')')
	w.Write(" VALUES ")
	i := 0
	for row := range rows {
		if row > 0 {
			w.Write(", ")
		}

		w.Char('(')
		for ci := range cols {
			if ci > 0 {
				w.Write(", ")
			}

			if !ph {
				w.Param([]any{})
			} else {
				w.Param([]any{result.Args[i]})
			}

			i++
		}
		w.Char(')')
	}
	return nil
}

// -- Methods

func (cls *InsertInto[Next]) Insert(reference ...any) Next {
	cls.declared = true
	cls.reference = reference
	return any(cls).(Next)
}

func (cls *InsertInto[Next]) Values(values any) {
	cls.values = values
}

func (cls *InsertInto[Next]) Table(table string) {
	cls.table = table
}
