package writer

import (
	"strings"

	"github.com/laacin/inyorm/intr/dialect"
)

type Dialect interface {
	dialect.ValueWriter
}

type Writer struct {
	sb                 strings.Builder
	ColumnWritingModes dialect.ClauseWritingConfig
	TableWritingModes  dialect.ClauseWritingConfig
	dial               Dialect
	ph                 *Placeholder
}

func (w *Writer) Write(str string) {
	w.sb.WriteString(str)
}

func (w *Writer) Char(char byte) {
	w.sb.WriteByte(char)
}

func (w *Writer) Value(v any, ctx dialect.ClauseName) {
	switch v := Normalize(v).(type) {
	case dialect.Param:
		w.ph.Store(v)
		w.sb.WriteString(w.dial.Param(w.ph.GetCount()))

	case dialect.Table:
		isDef := obtainWriteMode(ctx, &w.TableWritingModes) == dialect.WriteDef
		w.sb.WriteString(w.dial.Table(v, isDef))

	case dialect.Column:
		switch obtainWriteMode(ctx, &w.ColumnWritingModes) {
		case dialect.WriteBase:
			w.sb.WriteString(w.dial.ColBase(v))
		case dialect.WriteExpr:
			w.sb.WriteString(w.dial.ColExpr(v))
		case dialect.WriteAlias:
			w.sb.WriteString(w.dial.ColAlias(v))
		case dialect.WriteDef:
			w.sb.WriteString(w.dial.ColDef(v))
		default:
			w.sb.WriteString(w.dial.Null())
		}

	case string:
		w.sb.WriteString(w.dial.String(v))

	case int:
		w.sb.WriteString(w.dial.Number(v))

	case bool:
		w.sb.WriteString(w.dial.Bool(v))

	default:
		w.sb.WriteString(w.dial.Null())
	}
}

func (w *Writer) Result() string {
	return w.sb.String()
}

// -- Helpers
func obtainWriteMode(ctx dialect.ClauseName, modes *dialect.ClauseWritingConfig) dialect.WritingMode {
	switch ctx {
	case dialect.ClauseNameInsertInto:
		return modes.InsertInto

	case dialect.ClauseNameSelect:
		return modes.Select

	case dialect.ClauseNameFrom:
		return modes.From

	case dialect.ClauseNameJoin:
		return modes.Join

	case dialect.ClauseNameWhere:
		return modes.Where

	case dialect.ClauseNameGroupBy:
		return modes.GroupBy

	case dialect.ClauseNameHaving:
		return modes.Having

	case dialect.ClauseNameOrderBy:
		return modes.OrderBy

	case dialect.ClauseNameLimit:
		return modes.Limit

	case dialect.ClauseNameOffset:
		return modes.Offset

	case dialect.ClauseNameUpdate:
		return modes.Update

	case dialect.ClauseNameDelete:
		return modes.Delete

	default:
		return dialect.WriteBase
	}
}
