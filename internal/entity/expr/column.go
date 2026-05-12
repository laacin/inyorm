package expr

import "github.com/laacin/inyorm/internal/entity/core"

type Table struct{ Value string }

func (*Table) Kind() ValueKind { return ValueTable }

func (t *Table) Write(w core.InternalWriter, dial ValueSyntax, mode core.WritingMode) {
	w.SetRef(t.Value)
	dial.WriteTable(w, t)
}

type Column struct {
	Name  string // Column base name
	Ref   string // Table reference
	Alias string // Explicit alias
	Value string // Column expression
	From  Value
	Exprs []ColExpr
	Aggr  *ColExpr
}

func (c *Column) Kind() ValueKind { return ValueColumn }

func (c *Column) Write(w core.InternalWriter, dial ValueSyntax, mode core.WritingMode) {
	w.SetRef(c.Ref)

	switch mode {
	case core.WriteBase:
		dial.WriteColBase(w, c)
	case core.WriteExpr:
		dial.WriteColExpr(w, c)
	case core.WriteAlias:
		dial.WriteColAlias(w, c)
	case core.WriteDef:
		dial.WriteColDef(w, c)
	default:
		dial.WriteColExpr(w, c)
	}
}

// --- Tools for building
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

type ColExpr struct {
	Kind  ColKindExpr
	Value any // exists if is required. otherwise is nil
}
