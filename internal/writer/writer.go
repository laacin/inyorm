package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/core"
)

type Writer struct {
	sb        strings.Builder
	ph        *Placeholder
	aliases   *Alias
	colWriter *core.ColumnWriter
}

func (w *Writer) Write(v string) {
	w.sb.WriteString(v)
}

func (w *Writer) Char(v byte) {
	w.sb.WriteByte(v)
}

func (w *Writer) Param(v []any) {
	if len(v) > 0 {
		w.sb.WriteString(w.ph.withValue(v[0]))
		return
	}
	w.sb.WriteString(w.ph.write())
}

func (w *Writer) Value(v any, colWriter core.ColumnType) {
	switch val := v.(type) {
	case core.Builder:
		val(w)

	case core.Column:
		switch colWriter {
		case core.ColTypBase:
			val.Base()(w)

		case core.ColTypExpr:
			val.Expr()(w)

		case core.ColTypDef:
			val.Def()(w)

		case core.ColTypAlias:
			val.Alias()(w)
		}

	default:
		w.sb.WriteString(sqlify(val))
	}
}

func (w *Writer) Column(table, name string) {
	if w.aliases != nil && table != "" {
		w.sb.WriteByte(w.aliases.Get(table))
		w.sb.WriteByte('.')
	}
	w.sb.WriteString(name)
}

func (w *Writer) Table(v string) {
	w.sb.WriteString(v)
	if w.aliases != nil {
		w.sb.WriteByte(' ')
		w.sb.WriteByte(w.aliases.Get(v))
	}
}

func (w *Writer) Split() core.Writer {
	split := *w
	split.Reset()
	return &split
}

func (w *Writer) ToString() string {
	return w.sb.String()
}

func (w *Writer) Reset() {
	w.sb.Reset()
}
