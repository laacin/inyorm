package core

// ---- WRITER -----

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

type ValueParser interface {
	Render(InternalWriter, any, WritingMode)
}

// ---- PARAM STORE -----

type ParamStore interface {
	Store(v any)

	Lazy(id string)
	LazyObj(cols []string)

	Fill(id string, v any)
	FillObj(func(cols []string) []any) // Objs loads must be in orders

	LastIndex(idx int) ParamIndex
	Values() ([]any, error)
}

type ParamIndex struct {
	ID  string
	Num int
}

// ---- ALIAS STORE -----

type AliasStore interface {
	Enable()

	GetMain() Reference
	Get(ref string) Reference

	SetMain(ref string)
	Set(ref string)
}

type Reference struct {
	Ref     byte
	Enabled bool
}

// ---- MAPPER -----

type Mapper interface {
	ReadCols(...any) []string
	ReadValues(cols []string, v any) ([]any, error)
}
