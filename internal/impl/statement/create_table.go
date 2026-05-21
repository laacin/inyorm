package statement

import (
	"github.com/laacin/inyorm/internal/impl/table"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
)

type CreateTableQueryImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	table.TableBuilderImpl
}

func (s *CreateTableQueryImpl) Start(dial ir.Dialect, ref string) *CreateTableQueryImpl {
	s.DefaultRef = ref
	s.Dialect = dial
	s.TableBuilderImpl.Start(ref)
	return s
}

// --- Build

func (s *CreateTableQueryImpl) Build() (string, []any, error) {
	w := &writer.WriterImpl{Syntax: s.Dialect}
	s.TableBuilderImpl.Build(w, s.Dialect)
	return w.ToString(), nil, nil
}
