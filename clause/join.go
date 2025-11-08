package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

const (
	innerJoin = "INNER"
	leftJoin  = "LEFT"
	rightJoin = "RIGHT"
	fullJoin  = "FULL"
	crossJoin = "CROSS"
)

type JoinClause struct {
	joins []*OnJoin
}

func (j *JoinClause) Name() string {
	return core.ClsJoin
}

func (j *JoinClause) Build() core.Builder {
	return func(w core.Writer) {
		for i, join := range j.joins {
			if i > 0 {
				w.Char(' ')
			}
			w.Write(join.Type)
			w.Char(' ')
			w.Write("JOIN")
			w.Char(' ')
			w.Table(join.Table)
			if join.Cond != nil {
				w.Write(" ON ")
				join.Cond.Build(w, nil)
			}
		}
	}
}

// -- Methods

func (j *JoinClause) Join(table string) core.ClauseOn      { return j.join(innerJoin, table) }
func (j *JoinClause) JoinLeft(table string) core.ClauseOn  { return j.join(leftJoin, table) }
func (j *JoinClause) JoinRight(table string) core.ClauseOn { return j.join(rightJoin, table) }
func (j *JoinClause) JoinFull(table string) core.ClauseOn  { return j.join(fullJoin, table) }
func (j *JoinClause) JoinCross(table string)               { j.join(crossJoin, table) }

func (j *JoinClause) join(typ, table string) *OnJoin {
	join := &OnJoin{Type: typ, Table: table}
	j.joins = append(j.joins, join)
	return join
}

// -- Depending Clause

type OnJoin struct {
	Type  string
	Table string
	Cond  *condition.Condition
}

func (on *OnJoin) On(identifier any) core.Cond {
	cond := &condition.Condition{}
	on.Cond = cond
	return cond.Start(identifier)
}
