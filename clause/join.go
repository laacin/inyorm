package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

const (
	InnerJoin = "INNER"
	LeftJoin  = "LEFT"
	RightJoin = "RIGHT"
	FullJoin  = "FULL"
	CrossJoin = "CROSS"
)

type Join[Next, Cond, CondNext, Ident, Value any] struct {
	declared bool
	joins    []*join[Cond, CondNext, Ident, Value]
	current  *join[Cond, CondNext, Ident, Value]
}

func (cls *Join[Next, Cond, CondNext, Ident, Value]) Name() core.ClauseType { return core.ClsTypJoin }
func (cls *Join[Next, Cond, CondNext, Ident, Value]) IsDeclared() bool {
	return cls != nil && cls.declared
}
func (cls *Join[Next, Cond, CondNext, Ident, Value]) Build(w core.Writer) {
	for i, join := range cls.joins {
		if i > 0 {
			w.Char(' ')
		}
		w.Write(join.typ)
		w.Char(' ')
		w.Write("JOIN")
		w.Char(' ')
		w.Table(join.table)
		if join.cond != nil {
			w.Write(" ON ")
			join.cond.Build(w, cls.Name())
		}
	}
}

// -- Methods

func (cls *Join[Next, Cond, CondNext, Ident, Value]) Join(table string) Next {
	cls.declared = true
	join := &join[Cond, CondNext, Ident, Value]{typ: InnerJoin, table: table}
	cls.current = join
	cls.joins = append(cls.joins, join)
	return any(cls).(Next)
}

func (cls *Join[Next, Cond, CondNext, Ident, Value]) On(ident Ident) Cond {
	cond := &condition.Condition[Cond, CondNext, Ident, Value]{}
	cls.current.cond = cond
	return cond.Start(ident)
}

// -- internal

type join[Cond, CondNext, Ident, Value any] struct {
	typ   string
	table string
	cond  *condition.Condition[Cond, CondNext, Ident, Value]
}
