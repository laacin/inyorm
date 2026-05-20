package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/impl/table"
	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/ir"
)

type CreateTableStmtImpl struct {
	DefaultRef string
	Dialect    ir.Dialect

	table.TableBuilderImpl
}

func (s *CreateTableStmtImpl) Start(ctx context.Context, eng *ir.Engine, ref string) api.CreateTable {
	s.DefaultRef = ref
	s.Dialect = eng.Dialect
	s.TableBuilderImpl.Start(ref)
	return s
}

// --- Build

func (s *CreateTableStmtImpl) Build() string {
	w := &writer.WriterImpl{Syntax: s.Dialect}
	s.TableBuilderImpl.Build(w, s.Dialect)
	return w.ToString()
}
