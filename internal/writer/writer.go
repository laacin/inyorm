package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/core"
)

type Writer struct {
	sb      strings.Builder
	ph      *Placeholder
	aliases *Alias
}

func (w *Writer) Write(v string) {
	w.sb.WriteString(v)
}

func (w *Writer) Char(v byte) {
	w.sb.WriteByte(v)
}

func (w *Writer) Value(v any, opts *core.WriterOpts) {
	if opts.Placeholder {
		w.sb.WriteString(w.ph.next(v))
		return
	}

	switch val := v.(type) {
	case core.Builder:
		val(w)

	case core.Column:
		switch opts.ColType {
		case core.ColTypBase:
			val.Base()(w)
		case core.ColTypExpr:
			val.Expr()(w)
		case core.ColTypAlias:
			val.Alias()(w)
		case core.ColTypDef:
			val.Def()(w)
		}

	default:
		w.sb.WriteString(core.Sqlify(v))
	}
}

func (w *Writer) ColRef(table string) {
	ref := w.aliases.Get(table)
	w.sb.WriteByte(ref)
}

func (w *Writer) Table(v string) {
	w.sb.WriteString(v)
	w.sb.WriteByte(' ')
	w.sb.WriteByte(w.aliases.Get(v))
}

func (w *Writer) ToString() string {
	return w.sb.String()
}

func (w *Writer) Reset() {
	w.sb.Reset()
}
