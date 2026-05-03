package dialect

// Expr types
type (
	ExprOperator  int
	ExprConnector int
)

const (
	ExprEqual ExprOperator = iota
	ExprLike
	ExprIn
	ExprBetween
	ExprGreater
	ExprLess
	ExprIsNull

	ExprAnd ExprConnector = iota
	ExprOr
)

type Expr struct {
	Negated    bool
	Identifier any
	Operator   ExprOperator
	Values     []any
}

type Cond struct {
	Exprs      []Expr
	Connectors []ExprConnector
}
