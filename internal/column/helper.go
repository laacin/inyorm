package column

import "github.com/laacin/inyorm/internal/core"

func inferColumn[T any](w core.Writer, v any) {
	if col, ok := v.(*Column[T]); ok {
		col.Expr()(w)
		return
	}
	w.Identifier(v, core.ClsTypUnset)
}
