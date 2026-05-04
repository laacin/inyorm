package writer

import (
	"strings"

	"github.com/laacin/inyorm/intr/dialect"
)

type Writer struct {
	sb      strings.Builder
	aliases *AliasStore
	dial    dialect.Dialect
	params  *ParamStore
}

func (w *Writer) Write(str string) {
	w.sb.WriteString(str)
}

func (w *Writer) Char(char byte) {
	w.sb.WriteByte(char)
}

func (w *Writer) Value(v any, mode dialect.WritingMode) {
	switch v := Normalize(v).(type) {
	case dialect.Param:
		w.params.Store(v)
		ph := w.dial.Placeholder(w.params.GetCount())
		w.Write(ph)

	case dialect.Table:
		w.dial.Table(w, v)

	case dialect.Column:
		switch mode {
		case dialect.WriteBase:
			w.dial.ColWriteBase(w, v)

		case dialect.WriteExpr:
			w.dial.ColWriteExpr(w, v)

		case dialect.WriteAlias:
			w.dial.ColWriteAlias(w, v)

		case dialect.WriteDef:
			w.dial.ColWriteDef(w, v)

		default:
			w.dial.ColWriteExpr(w, v)
		}

	case string:
		r := w.dial.String(v)
		w.Write(r)

	case int:
		r := w.dial.Number(v)
		w.Write(r)

	case bool:
		r := w.dial.Bool(v)
		w.Write(r)

	default:
		r := w.dial.Null()
		w.Write(r)
	}
}

func (w *Writer) New() dialect.Writer {
	return &Writer{dial: w.dial}
}

func (w *Writer) GetTableRef(tbl string) (ref byte, shouldBeUsed bool) {
	if w.aliases == nil {
		return ' ', false
	}

	return w.aliases.Get(tbl), true
}

func (w *Writer) Result() string {
	return w.sb.String()
}

func (w *Writer) Reset() {
	w.sb.Reset()
}
