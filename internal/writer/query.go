package writer

import "github.com/laacin/inyorm/internal/core"

type Query struct {
	cfg          *core.Config
	table        string
	aliases      Alias
	placeholders Placeholder
	clauses      map[core.ClauseType]core.Clause
	clauseOrder  []core.ClauseType
}

func NewQuery(table string, cfg *core.Config) *Query {
	builder := &Query{
		table: table,
		cfg:   cfg,
	}
	builder.placeholders.dialect = cfg.Dialect

	if table != "" {
		builder.aliases.Get(table)
	}

	return builder
}

func (q *Query) Table() string { return q.table }
func (q *Query) Build() (string, []any) {
	w := &Writer{
		ph:        &q.placeholders,
		colWriter: &q.cfg.ColWrite,
	}

	if cls, exists := q.clauses[core.ClsTypJoin]; exists && cls.IsDeclared() {
		w.aliases = &q.aliases
	}

	i := 0
	for _, name := range q.clauseOrder {
		cls, exists := q.clauses[name]
		if !exists || !cls.IsDeclared() {
			continue
		}

		if i > 0 {
			w.Char(' ')
		}

		cls.Build(w, q.cfg)
		i++
	}

	return w.ToString(), q.placeholders.values
}

func (q *Query) SetClauses(clauses []core.Clause, order []core.ClauseType) {
	q.clauseOrder = order
	q.clauses = make(map[core.ClauseType]core.Clause, len(clauses))
	for _, cls := range clauses {
		q.clauses[cls.Name()] = cls
	}
}
