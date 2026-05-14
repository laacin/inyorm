package dml

import "github.com/laacin/inyorm/internal/core"

type Clause interface {
	Kind() ClauseKind
	Write(core.InternalWriter, ClauseSyntax)
}

// Wrapper implementations must implement this
type ClauseBuilder interface {
	IsDeclared() bool
	Kind() ClauseKind
	Build() (Clause, error)
}

type StatementBuilder interface {
	Kind() StatementKind
	Build() (*Statement, error)
}
