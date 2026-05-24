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
func (String) Kind() ExprKind { return ExprKindString }
func (Int) Kind() ExprKind    { return ExprKindNumber }
func (Float) Kind() ExprKind  { return ExprKindFloat }
func (Bool) Kind() ExprKind   { return ExprKindBool }
func (Null) Kind() ExprKind   { return ExprKindNull }

// Renders
func (v String) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteLitString(w, string(v))
}
func (v Int) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteLitInt(w, int(v))
}
func (v Float) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteLitFloat(w, float64(v))
}
func (v Bool) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteLitBool(w, bool(v))
}
func (Null) Render(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteLitNull(w)
}
