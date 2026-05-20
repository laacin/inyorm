package expr

type Table struct{ Value string }

type Column struct {
	Name  string // Column base name
	Ref   string // Table reference
	Alias string // Explicit alias
	Value string // Column expression
}

// --- Tools for building
type (
	ColExprKind int

	ColAggrKind   int
	ColScalarKind int
	ColArithKind  int
	ColWrapKind   int
)

const (
	// Kind
	ColKindAggr ColExprKind = iota
	ColKindScalar
	ColKindArith
	ColKindWrap

	// Aggregation
	ColAggrCount ColAggrKind = iota
	ColAggrSum
	ColAggrMin
	ColAggrMax
	ColAggrAvg

	// Scalar
	ColScalarLower ColScalarKind = iota
	ColScalarUpper
	ColScalarTrim
	ColScalarRound
	ColScalarAbs

	// Arith
	ColArithAdd ColArithKind = iota
	ColArithSub
	ColArithMul
	ColArithDiv
	ColArithMod
)

type ColExpr struct {
	Kind  ColExprKind
	Value any
}

type ColAggr struct {
	Kind     ColAggrKind
	Distinct bool
}

type ColScalar struct {
	Kind ColScalarKind
}

type ColArith struct {
	Kind  ColArithKind
	Value any
}

func (c *ColExpr) IsScalar() (*ColScalar, bool) {
	if c.Kind == ColKindScalar {
		v, ok := c.Value.(*ColScalar)
		return v, ok
	}
	return nil, false
}
func (c *ColExpr) IsArith() (*ColArith, bool) {
	if c.Kind == ColKindArith {
		v, ok := c.Value.(*ColArith)
		return v, ok
	}
	return nil, false
}
func (c *ColExpr) IsWrap() bool {
	return c.Kind == ColKindWrap
}
