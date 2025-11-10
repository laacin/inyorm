package clause

import "github.com/laacin/inyorm/internal/core"

type SelectClause struct {
	distinct bool
	targets  []any
}

func (s *SelectClause) Name() core.ClauseType {
	return core.ClsTypSelect
}

func (s *SelectClause) IsDeclared() bool { return s != nil }

func (s *SelectClause) Build(w core.Writer) {
	w.Write("SELECT ")

	if s.distinct {
		w.Write("DISTINCT ")
	}

	for i, sel := range s.targets {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(sel, core.WriterOpts{ColType: core.ColTypDef})
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
