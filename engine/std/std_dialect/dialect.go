package std_dialect

import (
	"github.com/laacin/inyorm/engine/std/std_dialect/std_dml"
	"github.com/laacin/inyorm/engine/std/std_dialect/std_expr"
)

type Dialect struct {
	std_expr.ExprSyntax
	std_dml.DmlSyntax
}
