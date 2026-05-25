package core

// Writing modes
type WritingMode int

const (
	WriteDef WritingMode = iota
	WriteBase
	WriteAlias
	WriteExpr
)

// Writer used by a dialect
type Writer interface {
	Write(string)
	Char(byte)
	Wrap(func(curr string, w Writer))
	Value(v any, mode WritingMode)

	GetRef(ref string) (byte, bool)
}

type InternalWriter interface {
	Writer
	SetRef(string)

	New() InternalWriter
	ToString() string
	Reset()
}
