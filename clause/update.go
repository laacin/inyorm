package clause

import (
	"errors"
	"fmt"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/mapper"
)

type Update[Next any] struct {
	declared  bool
	table     string
	reference []any
	values    any
}

func (cls *Update[Next]) Name() string     { return "UPDATE" }
func (cls *Update[Next]) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *Update[Next]) Build(w core.Writer, cfg *core.Config) error {
	if len(cls.reference) < 1 {
		return errors.New("missing reference")
	}

	cols, err := mapper.GetColumns(cfg.ColumnTag, cls.reference)
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

	ph := false
	if result != nil {
		if result.Rows > 1 {
			return errors.New("there should only be one row")
		}
		ph = result.Rows == 1
	}

	w.Write("UPDATE")
	w.Char(' ')
	w.Table(cls.table)
	w.Char(' ')
	w.Write("SET")
	w.Char(' ')
	for i, col := range cols {
		if i > 0 {
			w.Write(", ")
		}

		w.Column(cls.table, col)
		w.Write(" = ")
		if !ph {
			w.Param([]any{})
			continue
		}
		w.Param([]any{result.Args[i]})
	}
	return nil
}

// -- Methods

func (cls *Update[Next]) Update(reference ...any) Next {
	cls.declared = true
	cls.reference = reference
	return any(cls).(Next)
}

func (cls *Update[Next]) Values(values any) {
	cls.values = values
}

func (cls *Update[Next]) Table(table string) {
	cls.table = table
}
