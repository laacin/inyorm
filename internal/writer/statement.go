package writer

import "github.com/laacin/inyorm/internal/core"

type Statement struct {
	from    string
	dialect string
	ph      Placeholder
	alias   Alias
}

func NewStatement(dialect string) *Statement {
	return &Statement{
		dialect: dialect,
		ph:      Placeholder{dialect: dialect},
	}
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
