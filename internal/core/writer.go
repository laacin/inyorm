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
}

type InternalWriter interface {
	Writer
	New() InternalWriter
	ToString() string
	Reset()
}
