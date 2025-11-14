package writer

import "github.com/laacin/inyorm/internal/core"

type Statement struct {
	from    string
	dialect string
	ph      Placeholder
	alias   *Alias
	clauses map[core.ClauseType]core.Clause
}

func NewStatement(dialect string, defaultTable string) *Statement {
	stmt := &Statement{
		from:    defaultTable,
		dialect: dialect,
		ph:      Placeholder{dialect: dialect},
	}

	return stmt
}

func (stmt *Statement) Build(order []core.ClauseType) (string, []any) {
	if cls, exists := stmt.clauses[core.ClsTypJoin]; exists && cls.IsDeclared() {
		stmt.alias = &Alias{}
		stmt.alias.Get(stmt.from)
	}

	w := &Writer{ph: &stmt.ph, aliases: stmt.alias}
	i := 0
	for _, name := range order {
		cls, exists := stmt.clauses[name]
		if !exists || !cls.IsDeclared() {
			continue
		}

		if i > 0 {
			w.Char(' ')
		}
		cls.Build(w)
		i++
	}

	return w.ToString(), stmt.ph.values
}

func (stmt *Statement) SetClauses(clauses []core.Clause) {
	stmt.clauses = make(map[core.ClauseType]core.Clause, len(clauses))
	for _, cls := range clauses {
		stmt.clauses[cls.Name()] = cls
	}
}
