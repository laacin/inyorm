package clause

import "github.com/laacin/inyorm/internal/core"

type InsertInto struct {
	into   string
	cols   []any
	values []any
}

func (cls *InsertInto) Name() core.ClauseType { return core.ClsTypInsertInto }
func (cls *InsertInto) IsDeclared() bool      { return cls != nil }
func (cls *InsertInto) Build(w core.Writer) {
	var (
		colNum = len(cls.cols)
		valNum = len(cls.values)
	)
	if colNum < 1 || valNum%colNum != 0 {
		return
	}

	w.Write("INSERT INTO")
	w.Char(' ')
	w.Write(cls.into)
	w.Char(' ')
	w.Char('(')
	for i, col := range cls.cols {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(col, core.InsertIdentWriteOpt)
	}
	w.Char(')')
	w.Write(" VALUES ")

	rows := valNum / colNum
	i := 0
	for row := range rows {
		if row > 0 {
			w.Write(", ")
		}

		w.Char('(')
		for index := range colNum {
			if index > 0 {
				w.Write(", ")
			}
			w.Value(cls.values[i+index], core.InsertValueWriteOpt)
		}
		w.Char(')')
		i += colNum
	}
}

// -- Methods

func (cls *InsertInto) Insert(table string, cols []any) { cls.into = table; cls.cols = cols }
func (cls *InsertInto) Values(values []any)             { cls.values = values }
