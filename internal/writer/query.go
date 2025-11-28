package writer

import (
	"fmt"

	"github.com/laacin/inyorm/internal/core"
)

type Query struct {
	Config   *core.Config
	clauses  []core.Clause
	preBuild func(*core.Config) bool
}

func (q *Query) PreBuild(fn func(cfg *core.Config) (useAliases bool)) {
	q.preBuild = fn
}

func (q *Query) SetClauses(clauses []core.Clause) {
	q.clauses = clauses
}

func (q *Query) Build() (string, []any, error) {
	var (
		aliases *Alias
		phs     = &Placeholder{dialect: q.Config.Dialect}
	)

	if q.preBuild != nil && q.preBuild(q.Config) {
		aliases = &Alias{}
	}

	w := &Writer{
		colWriter: &q.Config.ColWrite,
		ph:        phs,
		aliases:   aliases,
	}

	i := 0
	for _, cls := range q.clauses {
		if !cls.IsDeclared() {
			continue
		}

		if i > 0 {
			w.Char(' ')
		}

		if err := cls.Build(w, q.Config); err != nil {
			return "", nil, fmt.Errorf("error building clause %s: %w", cls.Name(), err)
		}
		i++
	}

	return w.ToString(), phs.values, nil
}
