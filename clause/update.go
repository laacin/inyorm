package clause

import (
	"fmt"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/mapper"
)

type Update struct {
	declared bool
	table    string
	binder   any
}

func (cls *Update) Name() string     { return "UPDATE" }
func (cls *Update) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *Update) Build(w core.Writer, cfg *core.Config) error {
	result, err := mapper.Read(cfg.ColumnTag, cls.binder)
	if err != nil {
		return fmt.Errorf("failed to map value: %w", err)
	}

	var (
		cols = result.Columns
		vals = result.Args
	)

	upds := make([]update, len(cols))
	for i := range len(cols) {
		upds[i] = update{src: cols[i], value: vals[i]}
	}

	w.Write("UPDATE")
	w.Char(' ')
	w.Table(cls.table)
	w.Char(' ')
	w.Write("SET")
	w.Char(' ')
	for i, u := range upds {
		if i > 0 {
			w.Write(", ")
		}
		w.Column(cls.table, u.src)
		w.Write(" = ")
		// TODO: make dinamic param
		w.Param([]any{u.value})
		//w.Value(u.value, core.ColTypUnset)
	}
	return nil
}

// -- Methods

func (cls *Update) Update(binder any) {
	cls.declared = true
	cls.binder = binder
}

func (cls *Update) To(table string) {
	cls.table = table
}

// -- internal

type update struct {
	src   string
	value any
}
