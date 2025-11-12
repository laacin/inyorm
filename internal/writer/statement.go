package writer

import "github.com/laacin/inyorm/internal/core"

type Statement struct {
	from    string
	dialect string
	ph      Placeholder
	alias   Alias
	clauses []core.Clause
}

func NewStatement(dialect string, defaultTable string) *Statement {
	stmt := &Statement{
		from:    defaultTable,
		dialect: dialect,
		ph:      Placeholder{dialect: dialect},
	}
	stmt.alias.Get(defaultTable)

	return stmt
}

func (stmt *Statement) Build() (string, []any) {
	w := stmt.Writer()
	for i, cls := range stmt.clauses {
		if !cls.IsDeclared() {
			continue
		}

		if i > 0 {
			w.Char(' ')
		}
		cls.Build(w)
	}

	return w.ToString(), stmt.Values()
}

func (stmt *Statement) SetClauses(clauses []core.Clause) {
	stmt.clauses = clauses
}

func (stmt *Statement) Dialect() string { return stmt.dialect }

func (stmt *Statement) SetFrom(ref string) {
	stmt.from = ref
	stmt.alias.Get(ref)
}

func (stmt *Statement) GetFrom() string { return stmt.from }

func (stmt *Statement) Writer() core.Writer {
	return &Writer{ph: &stmt.ph, aliases: &stmt.alias}
}

func (stmt *Statement) Values() []any { return stmt.ph.values }
