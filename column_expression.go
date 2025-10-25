package inyorm

import (
	"strings"
	"sync"

	"github.com/laacin/inyorm/internal/stmt"
)

// ----- FIELD -----

// Represents a basic SQL field, used as an argument in queries.
type Field string

// Returns a valid SQL-safe value for use in a query.
func (f Field) Use() string { return stmt.SetColumn(string(f)) }

// ----- CUSTOM FIELD -----

type ColumnExpr struct{}

func (ce *ColumnExpr) New(alias string, fn func(fb *FB) *NdF) *CustomField {
	fb := &FB{}
	node := fn(fb)

	return &CustomField{
		sel:   node.build(),
		alias: alias,
	}
}

// Custom field created using a Field Builder (FB).
type CustomField struct {
	sel   string
	alias string
}

// Returns the SQL expression for this field, ready to be used in a SELECT clause.
func (n *CustomField) Select() string { return n.sel }

// Returns the alias of the column, properly formatted for use in a query.
func (n *CustomField) Alias() string { return stmt.SetColumn(n.alias) }

// ----- FIELD BUILDER -----

// Field Builder (FB) provides methods to define and build custom column expressions.
type FB struct {
	mu sync.Mutex
	sb strings.Builder
}

// Creates a node from a simple value.
// SQL: v
func (fb *FB) Simple(v Value) *NdF {
	target := &NdF{target: stmt.Stringify(v)}
	return target
}

// Creates a node from a concatenation of values.
// SQL: CONCAT(vals[0], vals[1], ...)
func (fb *FB) Concat(vals ...Value) *NdF {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	fb.sb.Reset()
	fb.sb.WriteString("CONCAT(")
	for i, v := range vals {
		if i > 0 {
			fb.sb.WriteString(", ")
		}
		writeValue(&fb.sb, v)
	}
	fb.sb.WriteByte(')')

	return &NdF{target: fb.sb.String()}
}

// Creates a new expression builder.
func (fb *FB) Expr(when Value) *Expr {
	expr := &Expr{ph: &stmt.PlaceholderGen{StringMode: true}}
	return expr.start(when)
}

// Represents a CASE field context used to define WHEN/THEN and ELSE clauses.
type CaseField[T any] struct {
	args []*Do[T]
	els  Value
}

// Defines a WHEN condition within a CASE clause.
// Returns a *Do[T] instance used to define the THEN action.
// SQL: WHEN v THEN ...
func (cs *CaseField[T]) When(v T) *Do[T] {
	do := &Do[T]{ctx: cs, when: v}
	cs.args = append(cs.args, do)
	return do
}

// Defines the ELSE branch of a CASE clause.
// SQL: ELSE v
func (cs *CaseField[T]) Else(v Value) {
	cs.els = v
}

// Represents a single WHEN/THEN pair within a CASE clause.
type Do[T any] struct {
	ctx  *CaseField[T]
	when T
	do   Value
}

// Defines the THEN action for a WHEN condition and returns the parent CaseField.
// Allows chaining multiple WHEN/THEN pairs.
// SQL: WHEN <when> THEN v
func (do *Do[T]) Then(v Value) *CaseField[T] {
	do.do = v
	return do.ctx
}

// Creates a node from a simple CASE statement.
// Starts from a condition and compares it against literal values.
// SQL: CASE cond WHEN ... THEN ... ELSE ... END
func (fb *FB) Switch(cond Value, fn func(cs *CaseField[Value])) *NdF {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	fb.sb.Reset()

	fb.sb.WriteString("CASE ")
	writeValue(&fb.sb, cond)
	cs := &CaseField[Value]{}

	fn(cs)
	for _, arg := range cs.args {
		fb.sb.WriteString(" WHEN ")
		writeValue(&fb.sb, arg.when)
		fb.sb.WriteString(" THEN ")
		writeValue(&fb.sb, arg.do)
	}
	fb.sb.WriteString(" ELSE ")
	writeValue(&fb.sb, cs.els)
	fb.sb.WriteString(" END")

	fld := &NdF{target: fb.sb.String()}
	return fld
}

// Creates a node from a searched CASE statement.
// Each WHEN clause is a full logical expression (built with fb.Expr()).
// SQL: CASE WHEN expr THEN ... ELSE ... END
func (fb *FB) Search(fn func(cs *CaseField[*ExprEnd])) *NdF {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	fb.sb.Reset()

	fb.sb.WriteString("CASE WHEN ")
	cs := &CaseField[*ExprEnd]{}
	fn(cs)
	for _, arg := range cs.args {
		arg.when.ctx.build(&fb.sb)
		fb.sb.WriteString(" THEN ")
		writeValue(&fb.sb, arg.do)
	}
	fb.sb.WriteString(" ELSE ")
	writeValue(&fb.sb, cs.els)
	fb.sb.WriteString(" END")

	fld := &NdF{target: fb.sb.String()}
	return fld
}

// ----- NODE FIELD -----

// NdF (Node field) is an internal type used to generate custom fields.
type NdF struct {
	target   string
	segments []nodeSegment
	aggr     string
}

type nodeSegment struct {
	typ   int
	arg   any
	value Value
}

// ----- AGGREGATION METHODS -----
const (
	countAggr = "COUNT"
	sumAggr   = "SUM"
	minAggr   = "MIN"
	maxAggr   = "MAX"
	avgAggr   = "AVG"
)

// Applies the COUNT aggregation to the current node.
// SQL: COUNT(<expr>)
func (n *NdF) Count() *NdF {
	n.aggr = countAggr
	return n
}

// Applies the SUM aggregation to the current node.
// SQL: SUM(<expr>)
func (n *NdF) Sum() *NdF {
	n.aggr = sumAggr
	return n
}

// Applies the MIN aggregation to the current node.
// SQL: MIN(<expr>)
func (n *NdF) Min() *NdF {
	n.aggr = minAggr
	return n
}

// Applies the MAX aggregation to the current node.
// SQL: MAX(<expr>)
func (n *NdF) Max() *NdF {
	n.aggr = maxAggr
	return n
}

// Applies the AVG aggregation to the current node.
// SQL: AVG(<expr>)
func (n *NdF) Avg() *NdF {
	n.aggr = avgAggr
	return n
}

const (
	scalarFn = iota
	arithOp
	wrap
)

// ----- ARITHMETICAL OPERATION METHODS -----
const (
	addOp byte = '+'
	subOp byte = '-'
	mulOp byte = '*'
	divOp byte = '/'
	modOp byte = '%'
)

// Adds a value to the current node using the addition operator (+).
// SQL: <expr> + v
func (n *NdF) Add(v Value) *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: arithOp, arg: addOp, value: v,
	})
	return n
}

// Subtracts a value from the current node using the subtraction operator (-).
// SQL: <expr> - v
func (n *NdF) Sub(v Value) *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: arithOp, arg: subOp, value: v,
	})
	return n
}

// Multiplies the current node by a given value using the multiplication operator (*).
// SQL: <expr> * v
func (n *NdF) Mul(v Value) *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: arithOp, arg: mulOp, value: v,
	})
	return n
}

// Divides the current node by a given value using the division operator (/).
// SQL: <expr> / v
func (n *NdF) Div(v Value) *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: arithOp, arg: divOp, value: v,
	})
	return n
}

// Applies the modulo operator (%) to the current node with the given value.
// SQL: <expr> % v
func (n *NdF) Mod(v Value) *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: arithOp, arg: modOp, value: v,
	})
	return n
}

// Wraps the current node in parentheses to control expression precedence.
// SQL: (<expr>)
func (n *NdF) Wrap() *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: wrap,
	})
	return n
}

// ----- SCALAR METHODS -----
const (
	lowerFunc = "LOWER"
	upperFunc = "UPPER"
	trimFunc  = "TRIM"
	roundFunc = "ROUND"
	absFunc   = "ABS"
)

// Converts the current expression to lowercase.
// SQL: LOWER(<expr>)
func (n *NdF) Lower() *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: scalarFn, arg: lowerFunc,
	})
	return n
}

// Converts the current expression to uppercase.
// SQL: UPPER(<expr>)
func (n *NdF) Upper() *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: scalarFn, arg: upperFunc,
	})
	return n
}

// Removes leading and trailing spaces from the current expression.
// SQL: TRIM(<expr>)
func (n *NdF) Trim() *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: scalarFn, arg: trimFunc,
	})
	return n
}

// Rounds a numeric expression to the nearest integer.
// SQL: ROUND(<expr>)
func (n *NdF) Round() *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: scalarFn, arg: roundFunc,
	})
	return n
}

// Returns the absolute value of a numeric expression.
// SQL: ABS(<expr>)
func (n *NdF) Abs() *NdF {
	n.segments = append(n.segments, nodeSegment{
		typ: scalarFn, arg: absFunc,
	})
	return n
}

// ----- Internal builders/helpers -----
func (n *NdF) build() string {
	var sb strings.Builder

	tgt := n.target
	for _, seg := range n.segments {
		sb.Reset()

		switch seg.typ {
		case scalarFn:
			sb.WriteString(seg.arg.(string))
			sb.WriteByte('(')
			sb.WriteString(tgt)
			sb.WriteByte(')')

		case arithOp:
			sb.WriteString(tgt)
			sb.WriteByte(' ')
			sb.WriteByte(seg.arg.(byte))
			sb.WriteByte(' ')
			writeValue(&sb, seg.value)

		case wrap:
			sb.WriteByte('(')
			sb.WriteString(tgt)
			sb.WriteByte(')')
		}
		tgt = sb.String()
	}

	if n.aggr != "" {
		sb.Reset()
		sb.WriteString(n.aggr)
		sb.WriteByte('(')
		sb.WriteString(tgt)
		sb.WriteByte(')')
		tgt = sb.String()
	}

	return tgt
}

func writeValue(sb *strings.Builder, v Value) {
	if n, ok := v.(*NdF); ok {
		sb.WriteString(n.build())
		return
	}
	sb.WriteString(stmt.Stringify(v))
}
