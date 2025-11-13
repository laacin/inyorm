package clause

import "github.com/laacin/inyorm/internal/core"

type Select struct {
	distinct bool
	targets  []any
}

func (s *Select) Name() core.ClauseType { return core.ClsTypSelect }
func (s *Select) IsDeclared() bool      { return s != nil }
func (s *Select) Build(w core.Writer) {
	w.Write("SELECT")
	w.Char(' ')

	if s.distinct {
		w.Write("DISTINCT ")
	}

	for i, sel := range s.targets {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(sel, core.SelectWriteOpt)
	}
}

// -- Methods

func (s *Select) Distinct() {
	s.distinct = true
}

func (s *Select) Select(targets []any) {
	s.targets = append(s.targets, targets...)
}
