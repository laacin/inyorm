package clause

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/mapper"
)

type InsertInto struct {
	declared bool
	table    string
	binder   any
}

func (cls *InsertInto) Name() core.ClauseType { return core.ClsTypInsert }
func (cls *InsertInto) IsDeclared() bool      { return cls != nil && cls.declared }

func (cls *InsertInto) Build(w core.Writer, cfg *core.Config) {
	// TODO: handle error
	result, _ := mapper.Read(cfg.ColumnTag, cls.binder)

	var (
		rows = result.Rows
		cols = result.Columns
		vals = result.Args
	)

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
		for _, val := range vals {
			if i > 0 {
				w.Write(", ")
			}
			w.Param([]any{val})
			i++
		}
		w.Char(')')
		i = 0
	}
}

// -- Methods

func (cls *InsertInto) Insert(binder any) {
	cls.declared = true
	cls.binder = binder
}

func (cls *InsertInto) Into(table string) {
	cls.table = table
}
