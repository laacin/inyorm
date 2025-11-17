package writer

import "github.com/laacin/inyorm/internal/core"

type StatementBuilder struct {
	Dialect      string
	Table        string
	Aliases      Alias
	Placeholders Placeholder
	Clauses      map[core.ClauseType]core.Clause
	ClauseOrder  []core.ClauseType
}

func NewStatement(dialect string, defaultTable string) *StatementBuilder {
	builder := &StatementBuilder{
		Table:   defaultTable,
		Dialect: dialect,
	}

	builder.Placeholders.dialect = dialect
	if defaultTable != "" {
		builder.Aliases.Get(defaultTable)
	}

	return builder
}

func (sb *StatementBuilder) Build() (string, []any) {
	w := &Writer{ph: &sb.Placeholders}

	if cls, exists := sb.Clauses[core.ClsTypJoin]; exists && cls.IsDeclared() {
		w.aliases = &sb.Aliases
	}

	i := 0
	for _, name := range sb.ClauseOrder {
		cls, exists := sb.Clauses[name]
		if !exists || !cls.IsDeclared() {
			continue
		}

		if i > 0 {
			w.Char(' ')
		}

		cls.Build(w)
		i++
	}

	return w.ToString(), sb.Placeholders.values
}

func (sb *StatementBuilder) SetClauses(clauses []core.Clause, order []core.ClauseType) {
	sb.ClauseOrder = order
	sb.Clauses = make(map[core.ClauseType]core.Clause, len(clauses))
	for _, cls := range clauses {
		sb.Clauses[cls.Name()] = cls
	}
}
