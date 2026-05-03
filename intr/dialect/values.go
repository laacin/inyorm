package dialect

// Simple values
type ValueKind int

// NOTE:: unused
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

// --- Types

// Param represents a value that will be stored and written as a dialect placeholder
type Param struct{ Value any }
