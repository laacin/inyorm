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

	var sb strings.Builder
	node.build(&sb)

	return &CustomField{
		sel:   sb.String(),
		alias: alias,
	}
}

// Custom field created using a Field Builder (FB).
type CustomField struct {
	sel   string
	alias string
}

// Returns the SQL expression for this field, ready to be used in a SELECT clause.
func (cf *CustomField) Select() string { return cf.sel }

// Returns the alias of the column, properly formatted for use in a query.
func (cf *CustomField) Alias() string { return stmt.SetColumn(cf.alias) }

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
	expr := &Expr{ph: &stmt.PlaceholderGen{Stringify: true}}
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
// Each WHEN clause inside the callback is an expression (created with fb.Expr()).
// SQL: CASE WHEN expr THEN ... ELSE ... END
func (fb *FB) Search(fn func(cs *CaseField[*ExprEnd])) *NdF {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	fb.sb.Reset()

	fb.sb.WriteString("CASE WHEN ")
	cs := &CaseField[*ExprEnd]{}
	fn(cs)
	for _, arg := range cs.args {
		arg.when.ctx.build(&fb.sb) // TODO: add support for Node
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

// Internal types
const (
	addOp  byte = '+'
	subOp  byte = '-'
	mulOp  byte = '*'
	divOp  byte = '/'
	modOp  byte = '%'
	wrapOp byte = 'W'
)

type operation struct {
	op    byte
	value any
}

// NdF (Node field) is an internal type used to generate custom fields.
type NdF struct {
	target string
	ops    []operation
	wraps  int
}

// Adds a value to the current node using the addition operator (+).
// SQL: <expr> + v
func (cf *NdF) Add(v Value) *NdF {
	cf.ops = append(cf.ops, operation{op: addOp, value: v})
	return cf
}

// Subtracts a value from the current node using the subtraction operator (-).
// SQL: <expr> - v
func (cf *NdF) Sub(v Value) *NdF {
	cf.ops = append(cf.ops, operation{op: subOp, value: v})
	return cf
}

// Multiplies the current node by a given value using the multiplication operator (*).
// SQL: <expr> * v
func (cf *NdF) Mul(v Value) *NdF {
	cf.ops = append(cf.ops, operation{op: mulOp, value: v})
	return cf
}

// Divides the current node by a given value using the division operator (/).
// SQL: <expr> / v
func (cf *NdF) Div(v Value) *NdF {
	cf.ops = append(cf.ops, operation{op: divOp, value: v})
	return cf
}

// Applies the modulo operator (%) to the current node with the given value.
// SQL: <expr> % v
func (cf *NdF) Mod(v Value) *NdF {
	cf.ops = append(cf.ops, operation{op: modOp, value: v})
	return cf
}

// Wraps the current node in parentheses to control expression precedence.
// SQL: (<expr>)
func (cf *NdF) Wrap() *NdF {
	cf.wraps++
	cf.ops = append(cf.ops, operation{op: wrapOp})
	return cf
}

// ---- Internal builders/helpers -----

func (cf *NdF) build(sb *strings.Builder) {
	cf.buildOperation(sb)
}

func (cf *NdF) buildOperation(sb *strings.Builder) {
	if cf.ops == nil {
		sb.WriteString(cf.target)
		return
	}

	for range cf.wraps {
		sb.WriteByte('(')
	}

	sb.WriteByte('(')
	sb.WriteString(cf.target)
	sb.WriteByte(')')

	for _, v := range cf.ops {
		if v.op == wrapOp {
			sb.WriteByte(')')
			continue
		}

		sb.WriteByte(' ')
		sb.WriteByte(v.op)
		sb.WriteByte(' ')
		sb.WriteString(stmt.Stringify(v.value))
	}
}

func writeValue(sb *strings.Builder, v Value) {
	if cf, ok := v.(*NdF); ok {
		cf.build(sb)
		return
	}
	sb.WriteString(stmt.Stringify(v))
}
