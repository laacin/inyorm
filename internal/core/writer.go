package core

type Writer interface {
	Write(v string)
	Char(v byte)

	Value(v any, opts ValueOpts)
	ColRef(table string)
	Table(v string)

	ToString() string
	Reset()
}

type ValueOpts struct {
	Placeholder bool
	Definition  bool
}
