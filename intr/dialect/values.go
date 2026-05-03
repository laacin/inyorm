package dialect

// Simple values
type ValueKind int

const (
	ValueString ValueKind = iota
	ValueNumber
	ValueFloat
	ValueBool
	ValueNull
	ValueParam

	ValueTable
	ValueColumn
	ValueExpr
)

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

// Writing form
type WritingMode int

const (
	WriteDef WritingMode = iota
	WriteBase
	WriteAlias // Column only
	WriteExpr  // Column only
)

var WritingFallbacks = [...]WritingMode{
	WriteDef,
	WriteAlias,
	WriteExpr,
	WriteBase,
}

// --- SQL Types

type Param struct {
	Value any // will be saved and write dialect parameter instead
}

type Table struct {
	Name string // Table name
	Ref  byte   // Alias reference (Table ref)
}

type Column struct {
	Name  string // Column base name
	Ref   string // Table reference
	Alias string // Explicit alias
	Value string // Column expression
}

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
