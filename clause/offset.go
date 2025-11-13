package clause

import "github.com/laacin/inyorm/internal/core"

type Offset struct {
	offset int
}

func (o *Offset) Name() core.ClauseType { return core.ClsTypOffset }
func (o *Offset) IsDeclared() bool      { return o != nil }
func (o *Offset) Build(w core.Writer) {
	if o.offset > 0 {
		w.Write("OFFSET")
		w.Char(' ')
		w.Value(o.offset, core.OffsetWriteOpt)
	}
}

// -- Methods

func (o *Offset) Offset(value int) {
	if value > 0 {
		o.offset = value
	}
}
