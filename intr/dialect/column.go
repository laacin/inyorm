package dialect

// --- Main types
type Table struct{ Name string }

type Column struct {
	Name  string // Column base name
	Table string // Table reference
	Alias string // Explicit alias
	Value string // Column expression
}

type ColKindExpr int

const (
	// Aggregation
	ColAggrCount ColKindExpr = iota
	ColAggrSum
	ColAggrMin
	ColAggrMax
	ColAggrAvg

	// Scalar
	ColScalarLower
	ColScalarUpper
	ColScalarTrim
	ColScalarRound
	ColScalarAbs

	// Arith
	ColArithAdd
	ColArithSub
	ColArithMul
	ColArithDiv
	ColArithMod
	ColArithWrap
)

// --- Column essentials

type ColExpr struct {
	Current string
	Kind    ColKindExpr
	Value   any // exists if is required. otherwise is nil
}

type CaseExpr struct {
	Identifier any
	Argument   any
}

type CaseCond struct {
	Exprs []CaseExpr
	Els   any
}

type ColumnEssentials interface {
	EssColExpr(Writer, ColExpr) (string, error)
}
