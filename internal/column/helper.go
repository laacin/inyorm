package column

import "github.com/laacin/inyorm/internal/core"

func inferColumn[T any](w core.Writer, v any) {
	if col, ok := v.(*Column[T]); ok {
		col.Expr()(w)
		return
	}
	w.Value(v, core.ClsTypUnset)
}

func tbl(dflt string, provided []string) string {
	if len(provided) > 0 {
		return provided[0]
	}
	return dflt
}
