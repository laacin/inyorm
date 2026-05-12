package dml

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

type ClauseKind int

const (
	// Select statement
	ClauseSelect ClauseKind = iota
	ClauseFrom
	ClauseJoin
	ClauseWhere
	ClauseGroupBy
	ClauseHaving
	ClauseOrderBy
	ClauseLimit
	ClauseOffset

	// Insert statement
	ClauseInsertInto

	// Update statement
	ClauseUpdate

	// Delete statement
	ClauseDelete
)

type StatementKind int

const (
	StatementSelect StatementKind = iota
	StatementInsert
	StatementUpdate
	StatementDelete
)

type Value interface {
	Kind() ValueKind
	Write(core.InternalWriter, ValueSyntax, core.WritingMode)
}

type Clause interface {
	Kind() ClauseKind
	Write(core.InternalWriter, ClauseSyntax)
}

// Wrapper implementations must implement this
type ValueBuilder interface {
	Build() Value
}

type ClauseBuilder interface {
	IsDeclared() bool
	Kind() ClauseKind
	Build() (Clause, error)
}

type StatementBuilder interface {
	Kind() StatementKind
	Build() (*Statement, error)
}
