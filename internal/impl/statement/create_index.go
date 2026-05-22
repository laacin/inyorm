package statement

import (
	"github.com/laacin/inyorm/bkp/table"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
)

type CreateIndexQueryImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	table.ConsDeclImpl
}

func (s *CreateIndexQueryImpl) Start(dial ir.Dialect, ref string) *CreateIndexQueryImpl {
	s.DefaultRef = ref
	s.Dialect = dial
	return s
}

// --- Build

func (s *CreateIndexQueryImpl) Build() string {
	w := &writer.WriterImpl{Syntax: s.Dialect}

	// s.TableBuilderImpl.Build(w, s.Dialect)
	return w.ToString()
}
