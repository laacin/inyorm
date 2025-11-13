package core

type Builder = func(Writer)

type Writer interface {
	Write(v string)
	Char(v byte)

	Value(v any, opts *WriterOpts)
	ColRef(table string)
	Table(v string)

	ToString() string
	Reset()
}

type WriterOpts struct {
	Placeholder bool
	ColType     ColumnType
}

// Column writer options
var (
	ColumnIdentWriteOpt = &WriterOpts{ColType: ColTypExpr}
	ColumnValueWriteOpt = &WriterOpts{}
)

// Clause writer options
var (
	InsertIdentWriteOpt = &WriterOpts{ColType: ColTypBase}
	InsertValueWriteOpt = &WriterOpts{Placeholder: true}

	SelectWriteOpt      = &WriterOpts{ColType: ColTypDef}
	JoinIdentWriteOpt   = &WriterOpts{ColType: ColTypBase}
	JoinValueWriteOpt   = &WriterOpts{}
	WhereIdentWriteOpt  = &WriterOpts{ColType: ColTypExpr}
	WhereValueWriteOpt  = &WriterOpts{Placeholder: true}
	GroupByWriteOpt     = &WriterOpts{ColType: ColTypExpr}
	HavingIdentWriteOpt = &WriterOpts{ColType: ColTypExpr}
	HavingValueWriteOpt = &WriterOpts{}
	OrderByWriteOpt     = &WriterOpts{ColType: ColTypAlias}
	LimitWriteOpt       = &WriterOpts{}
	OffsetWriteOpt      = &WriterOpts{}
)
