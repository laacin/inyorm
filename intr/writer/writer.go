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
		w.dial.Placeholder(w, w.params.GetCount())

	case dialect.Table:
		w.dial.Table(w, v)

	case dialect.Column:
		switch mode {
		case dialect.WriteBase:
			w.dial.ColBase(w, v)

		case dialect.WriteExpr:
			w.dial.ColExpr(w, v)

		case dialect.WriteAlias:
			w.dial.ColAlias(w, v)

		case dialect.WriteDef:
			w.dial.ColDef(w, v)

		default:
			w.dial.ColExpr(w, v)
		}

	case string:
		w.dial.String(w, v)

	case int:
		w.dial.Number(w, v)

	case bool:
		w.dial.Bool(w, v)

	default:
		w.dial.Null(w)
	}
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
