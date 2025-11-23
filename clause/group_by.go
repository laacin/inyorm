package clause

import "github.com/laacin/inyorm/internal/core"

type GroupBy struct {
	declared bool
	groups   []any
}

func (cls *GroupBy) Name() core.ClauseType { return core.ClsTypGroupBy }
func (cls *GroupBy) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *GroupBy) Build(w core.Writer, cfg *core.Config) {
	w.Write("GROUP BY")
	w.Char(' ')
	for i, group := range cls.groups {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(group, cfg.ColWrite.GroupBy)
	}
}

// -- Methods

func (cls *GroupBy) GroupBy(value ...any) {
	cls.declared = true
	cls.groups = value
}
