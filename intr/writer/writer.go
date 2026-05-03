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
	params             *ParamStore
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
		w.params.Store(v)
		w.dial.Placeholder(w, w.params.GetCount())

	case dialect.Table:
		isDef := obtainWriteMode(ctx, &w.TableWritingModes) == dialect.WriteDef
		w.dial.Table(w, v, isDef)

	case dialect.Column:
		switch obtainWriteMode(ctx, &w.ColumnWritingModes) {
		case dialect.WriteBase:
			w.dial.ColBase(w, v)
		case dialect.WriteExpr:
			w.dial.ColExpr(w, v)
		case dialect.WriteAlias:
			w.dial.ColAlias(w, v)
		case dialect.WriteDef:
			w.dial.ColDef(w, v)
		default:
			w.dial.Null(w)
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
