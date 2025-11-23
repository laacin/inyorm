package clause

import "github.com/laacin/inyorm/internal/core"

type Limit struct {
	declared bool
	limit    int
}

func (cls *Limit) Name() core.ClauseType { return core.ClsTypLimit }
func (cls *Limit) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Limit) Build(w core.Writer) {
	w.Write("LIMIT")
	w.Char(' ')
	w.Value(cls.limit, cls.Name())
}

// -- Methods

func (cls *Limit) Limit(value int) {
	if value > 0 {
		cls.declared = true
		cls.limit = value
	}
}
