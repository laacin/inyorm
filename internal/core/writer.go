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
	Wrap(func(current string, w Writer))
	Value(v any, mode WritingMode)

	GetRef(ref string) (byte, bool)
	ValueCount() int
}

type InternalWriter interface {
	Writer
	PushValue(v any)
	PushLazyValue(ref string)
	PushLazyObj(cols []string)
	SetRef(string)

	New() InternalWriter
	ToString() string
	Reset()
}

type WriterFunc = func(Writer)
