package expr

import "github.com/laacin/inyorm/internal/entity/core"

type ValueKind int

const (
	// Literals
	ValueString ValueKind = iota
	ValueNumber
	ValueFloat
	ValueBool
	ValueNull

	// Specials
	ValueParameter
	ValueWildcard
	ValueCondition
	ValueConcat
	ValueCaseSwitch
	ValueCaseSearch

	// SQL Values
	ValueTable
	ValueColumn
)

type Value interface {
	Kind() ValueKind
	Write(core.InternalWriter, ValueSyntax, core.WritingMode)
}

// Wrapper implementations must implement this
type ValueBuilder interface {
	Build() Value
}
