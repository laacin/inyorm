package expr

import "github.com/laacin/inyorm/internal/core"

type (
	String string
	Number int
	Float  float64
	Bool   bool
	Null   struct{}
)

// Kinds
func (String) Kind() ExprKind { return ExprString }
func (Number) Kind() ExprKind { return ExprNumber }
func (Float) Kind() ExprKind  { return ExprFloat }
func (Bool) Kind() ExprKind   { return ExprBool }
func (Null) Kind() ExprKind   { return ExprNull }

// Writers
func (v String) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteString(w, string(v))
}
func (v Number) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteNumber(w, int(v))
}
func (v Float) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteFloat(w, float64(v))
}
func (v Bool) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteBool(w, bool(v))
}
func (Null) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	dial.WriteNull(w)
}
