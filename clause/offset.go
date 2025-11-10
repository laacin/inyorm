package clause

import "github.com/laacin/inyorm/internal/core"

type OffsetClause struct {
	offset int
}

func (o *OffsetClause) Name() core.ClauseType {
	return core.ClsTypOffset
}

func (o *OffsetClause) IsDeclared() bool { return o != nil }

func (o *OffsetClause) Build(w core.Writer) {
	if o.offset > 0 {
		w.Write("OFFSET ")
		w.Value(o.offset, core.WriterOpts{})
	}
}

// -- Methods

func (o *OffsetClause) Offset(value int) {
	if value > 0 {
		o.offset = value
	}
}
