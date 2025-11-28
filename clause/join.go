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

type Join[Next, End, Cond, CondNext any] struct {
	declared bool
	joins    []*join[Cond, CondNext]
	current  *join[Cond, CondNext]
}

func (cls *Join[Next, End, Cond, CondNext]) Name() string { return "JOIN" }
func (cls *Join[Next, End, Cond, CondNext]) IsDeclared() bool {
	return cls != nil && cls.declared
}
func (cls *Join[Next, End, Cond, CondNext]) Build(w core.Writer, cfg *core.Config) error {
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
			join.cond.Build(w, cfg.ColWrite.Join)
		}
	}
	return nil
}

// -- Methods

func (cls *Join[Next, End, Cond, CondNext]) Join(table string) Next {
	cls.declared = true
	join := &join[Cond, CondNext]{typ: InnerJoin, table: table}
	cls.joins = append(cls.joins, join)
	cls.current = join
	return any(cls).(Next)
}

func (cls *Join[Next, End, Cond, CondNext]) Left() End {
	cls.current.typ = LeftJoin
	return any(cls).(End)
}

func (cls *Join[Next, End, Cond, CondNext]) Right() End {
	cls.current.typ = RightJoin
	return any(cls).(End)
}

func (cls *Join[Next, End, Cond, CondNext]) Full() End {
	cls.current.typ = FullJoin
	return any(cls).(End)
}

func (cls *Join[Next, End, Cond, CondNext]) Cross() {
	cls.current.typ = CrossJoin
}

func (cls *Join[Next, End, Cond, CondNext]) On(ident any) Cond {
	cond := &condition.Condition[Cond, CondNext]{}
	cls.current.cond = cond
	return cls.current.cond.Start(ident)
}

// -- internal

type join[Cond, CondNext any] struct {
	typ   string
	table string
	cond  *condition.Condition[Cond, CondNext]
}
