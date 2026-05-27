package expr

import "github.com/laacin/inyorm/internal/core"

type (
	String string
	Int    int
	Float  float64
	Bool   bool
	Null   struct{}
)

// Kinds
func (String) Kind() Kind { return KindString }
func (Int) Kind() Kind    { return KindNumber }
func (Float) Kind() Kind  { return KindFloat }
func (Bool) Kind() Kind   { return KindBool }
func (Null) Kind() Kind   { return KindNull }

// Renders
func (v String) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteLitString(w, string(v))
}
func (v Int) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteLitInt(w, int(v))
}
func (v Float) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteLitFloat(w, float64(v))
}
func (v Bool) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteLitBool(w, bool(v))
}
func (Null) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteLitNull(w)
}
