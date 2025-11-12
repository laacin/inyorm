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

type Join struct {
	joins   []*join
	current *join
}

func (j *Join) Name() core.ClauseType { return core.ClsTypJoin }
func (j *Join) IsDeclared() bool      { return j != nil }
func (j *Join) Build(w core.Writer) {
	for i, join := range j.joins {
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
			join.cond.Build(w, core.WriterOpts{ColType: core.ColTypBase})
		}
	}
}

// -- Methods

func (j *Join) Join(typ, table string) {
	join := &join{typ: typ, table: table}
	j.current = join
	j.joins = append(j.joins, join)
}

func (j *Join) On(ident any) core.Condition {
	cond := &condition.Condition{}
	j.current.cond = cond
	cond.Start(ident)
	return cond
}

// -- internal

type join struct {
	typ   string
	table string
	cond  *condition.Condition
}
