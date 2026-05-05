package entity

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
	Write(Writer, ValueWriter, WritingMode)
}

// Wrapper implemetations must implements this
type ValueDefer interface {
	Defer() Value
}
