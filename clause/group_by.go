package clause

import "github.com/laacin/inyorm/internal/core"

type GroupBy struct {
	groups []any
}

func (g *GroupBy) Name() core.ClauseType { return core.ClsTypGroupBy }
func (g *GroupBy) IsDeclared() bool      { return g != nil }
func (g *GroupBy) Build(w core.Writer) {
	w.Write("GROUP BY")
	w.Char(' ')
	for i, group := range g.groups {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, core.GroupByWriteOpt)
	}
}

// -- Methods

func (g *GroupBy) GroupBy(value []any) {
	g.groups = value
}
