package clause

import "github.com/laacin/inyorm/internal/core"

type Offset struct {
	declared bool
	offset   int
}

func (cls *Offset) Name() core.ClauseType { return core.ClsTypOffset }
func (cls *Offset) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *Offset) Build(w core.Writer) {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(cls.offset, cls.Name())
}

// -- Methods

func (cls *Offset) Offset(value int) {
	if value > 0 {
		cls.declared = true
		cls.offset = value
	}
}
