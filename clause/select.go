package clause

import "github.com/laacin/inyorm/internal/core"

type SelectClause struct {
	distinct bool
	targets  []any
}

func (s *SelectClause) Name() string {
	return core.ClsSelect
}

func (s *SelectClause) Build() core.Builder {
	return func(w core.Writer) {
		w.Write("SELECT ")

		if s.distinct {
			w.Write("DISTINCT ")
		}

		for i, sel := range s.targets {
			if i > 0 {
				w.Write(", ")
			}
			w.Value(sel, &core.ValueOpts{Definition: true})
		}
	}
}

// -- Methods

func (s *SelectClause) Distinct() core.ClauseSelect {
	s.distinct = true
	return s
}

func (s *SelectClause) Select(targets ...any) {
	s.targets = append(s.targets, targets...)
}
