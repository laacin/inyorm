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

// --- SQL Types

type Param struct {
	Value any // will be saved and write dialect parameter instead
}
