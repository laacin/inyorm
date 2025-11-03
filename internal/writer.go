package internal

import "strings"

type Statement struct {
	writer  *StmtWriter
	clauses []func(*StmtWriter)
	order   []string
}

func NewStatement() *Statement {
	return &Statement{writer: &StmtWriter{}}
}

func (stmt *Statement) Clause(fn func(w *StmtWriter)) {
	stmt.clauses = append(stmt.clauses, fn)
}

func (stmt *Statement) Build(aliases *Alias, hasJoins bool) string {
	writer := &StmtWriter{
		aliases: aliases,
		hasJoin: hasJoins,
	}

	for i, cls := range stmt.clauses {
		if i > 0 {
			writer.Char(' ')
		}
		cls(writer)
	}

	return writer.sb.String()
}

type StmtWriter struct {
	sb      strings.Builder
	aliases *Alias
	hasJoin bool
}

func (w *StmtWriter) Write(v string) { w.sb.WriteString(v) }
func (w *StmtWriter) Char(v byte)    { w.sb.WriteByte(v) }
func (w *StmtWriter) InferValue(v any, def bool) {
	switch val := v.(type) {
	case Column:
		w.Column(val, def)
	default:
		w.Sqlify(val)
	}
}

func (w *StmtWriter) Sqlify(v any) { w.sb.WriteString(Sqlify(v)) }
func (w *StmtWriter) Column(v Column, definition bool) {
	if v.Custom {
		if definition {
			w.sb.WriteString(v.Value)
			w.sb.WriteString(" AS ")
		}
		w.sb.WriteString(v.Alias)
		return
	}

	if w.hasJoin {
		alias := w.aliases.Get(v.Table)
		w.sb.WriteByte(alias)
		w.sb.WriteByte('.')
	}
	w.sb.WriteString(v.Value)
}

func (w *StmtWriter) Table(v string) {
	w.sb.WriteString(v)
	if w.hasJoin {
		w.sb.WriteByte(' ')
		w.sb.WriteByte(w.aliases.Get(v))
	}
}
