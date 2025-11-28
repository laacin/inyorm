package clause

import "github.com/laacin/inyorm/internal/core"

type Offset struct {
	declared bool
	offset   int
}

func (cls *Offset) Name() string     { return "OFFSET" }
func (cls *Offset) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *Offset) Build(w core.Writer, cfg *core.Config) error {
	w.Write("OFFSET")
	w.Char(' ')
	w.Value(cls.offset, core.ColTypUnset)
	return nil
}

// -- Methods

func (cls *Offset) Offset(value int) {
	if value > 0 {
		cls.declared = true
		cls.offset = value
	}
}
