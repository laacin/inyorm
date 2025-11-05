package internal

const (
	psql = "postgres"
)

type Statement struct {
	tbl         string
	clauses     []func(*Writer)
	aliases     Alias
	placeholder Placeholder
}

type OptionStatement struct {
	Dialect      string
	DefaultTable string
}

func NewStatement(opts *OptionStatement) *Statement {
	stmt := &Statement{}
	if opts != nil {
		if tbl := opts.DefaultTable; tbl != "" {
			stmt.tbl = tbl
			stmt.aliases.Get(tbl)
		}

		if opts.Dialect != "" {
			stmt.placeholder.dialect = opts.Dialect
		}
	}
	return stmt
}

func (stmt *Statement) NewWriter() *Writer {
	return &Writer{
		ph:      &stmt.placeholder,
		aliases: &stmt.aliases,
	}
}

func (stmt *Statement) MainTable() string {
	return stmt.tbl
}

func (stmt *Statement) Clause(fn func(w *Writer)) {
	stmt.clauses = append(stmt.clauses, fn)
}

func (stmt *Statement) Build() (string, []any) {
	w := stmt.NewWriter()
	for i, cls := range stmt.clauses {
		if i > 0 {
			w.Char(' ')
		}
		cls(w)
	}

	return w.sb.String(), stmt.placeholder.values
}
