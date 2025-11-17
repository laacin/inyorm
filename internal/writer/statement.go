package writer

import "github.com/laacin/inyorm/internal/core"

type StatementBuilder struct {
	dialect      string
	table        string
	aliases      Alias
	placeholders Placeholder
	clauses      map[core.ClauseType]core.Clause
	clauseOrder  []core.ClauseType
}

func NewStatement(dialect string, defaultTable string) *StatementBuilder {
	builder := &StatementBuilder{
		table:   defaultTable,
		dialect: dialect,
	}

	builder.placeholders.dialect = dialect
	if defaultTable != "" {
		builder.aliases.Get(defaultTable)
	}

	return builder
}

func (sb *StatementBuilder) Build() (string, []any) {
	w := &Writer{ph: &sb.placeholders}

	if cls, exists := sb.clauses[core.ClsTypJoin]; exists && cls.IsDeclared() {
		w.aliases = &sb.aliases
	}

	i := 0
	for _, name := range sb.clauseOrder {
		cls, exists := sb.clauses[name]
		if !exists || !cls.IsDeclared() {
			continue
		}

		if i > 0 {
			w.Char(' ')
		}

		cls.Build(w)
		i++
	}

	return w.ToString(), sb.placeholders.values
}

func (sb *StatementBuilder) SetClauses(clauses []core.Clause, order []core.ClauseType) {
	sb.clauseOrder = order
	sb.clauses = make(map[core.ClauseType]core.Clause, len(clauses))
	for _, cls := range clauses {
		sb.clauses[cls.Name()] = cls
	}
}
