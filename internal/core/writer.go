package core

type Builder = func(Writer)

type Writer interface {
	Write(v string)
	Char(v byte)

	Value(v any, opts WriterOpts)
	ColRef(table string)
	Table(v string)

	ToString() string
	Reset()
}

type WriterOpts struct {
	Placeholder bool
	ColType     ColumnType
}
