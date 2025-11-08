package clause

import (
	"github.com/laacin/inyorm/internal/condition"
	"github.com/laacin/inyorm/internal/core"
)

type GroupByClause struct {
	groups []any
	having *HavingClause
}

func (g *GroupByClause) Name() string {
	return core.ClsGroupBy
}

func (g *GroupByClause) Build() core.Builder {
	return func(w core.Writer) {
		w.Write("GROUP BY")
		w.Char(' ')
		for i, group := range g.groups {
			if i > 0 {
				w.Write(", ")
			}
			w.Value(group, nil)
		}
		if cond := g.having.cond; cond != nil {
			w.Write(" HAVING ")
			cond.Build(w, &core.ValueOpts{Definition: true})
		}
	}
}

// -- Methods

func (g *GroupByClause) GroupBy(value ...any) core.ClauseHaving {
	g.groups = append(g.groups, value...)
	g.having = &HavingClause{}
	return g.having
}

// -- Depending clauses

type HavingClause struct {
	cond *condition.Condition
}

func (g *HavingClause) Having(cond any) core.Cond {
	g.cond = &condition.Condition{}
	return g.cond.Start(cond)
}
