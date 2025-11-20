package column

import "github.com/laacin/inyorm/internal/core"

func inferColumn[T, K any](w core.Writer, v any) {
	if col, ok := v.(*Column[T, K]); ok {
		col.Expr()(w)
		return
	}
	w.Identifier(v, core.ClsTypUnset)
}
