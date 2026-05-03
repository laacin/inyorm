package dialect

// --- Main types
type Table struct {
	Name string // Table name
	Ref  byte   // Alias reference (Table ref)
}

type Column struct {
	Name  string // Column base name
	Ref   *byte  // Table reference
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

type ColumnEssentials interface {
	BuildColExpr(Writer, ColExpr) (string, error)
}
