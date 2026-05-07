package entity

type (
	String   string
	Number   int
	Float    float64
	Bool     bool
	Null     struct{}
	Wildcard struct{}
)

// Kinds
func (String) Kind() ValueKind   { return ValueString }
func (Number) Kind() ValueKind   { return ValueNumber }
func (Float) Kind() ValueKind    { return ValueFloat }
func (Bool) Kind() ValueKind     { return ValueBool }
func (Null) Kind() ValueKind     { return ValueNull }
func (Wildcard) Kind() ValueKind { return ValueWildcard }

// Writers
func (v String) Write(w InternalWriter, dial ValueSyntax, mode WritingMode) {
	dial.WriteString(w, string(v))
}
func (v Number) Write(w InternalWriter, dial ValueSyntax, mode WritingMode) {
	dial.WriteNumber(w, int(v))
}
func (v Float) Write(w InternalWriter, dial ValueSyntax, mode WritingMode) {
	dial.WriteFloat(w, float64(v))
}
func (v Bool) Write(w InternalWriter, dial ValueSyntax, mode WritingMode) {
	dial.WriteBool(w, bool(v))
}
func (Null) Write(w InternalWriter, dial ValueSyntax, mode WritingMode) {
	dial.WriteNull(w)
}
func (v *Wildcard) Write(w InternalWriter, dial ValueSyntax, mode WritingMode) {
	dial.WriteWildcard(w)
}
