package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal"
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/writer"
)

type selectStmt struct {
	stmt *internal.Statement
	dial entity.Dialect

	clsSelect
	clsFrom
	clsJoin
	clsWhere
	clsGroupBy
	clsHaving
	clsOrderBy
	clsLimit
	clsOffset
}

func newSelect(ctx context.Context, dial entity.Dialect, table string) *selectStmt {
	stmt := internal.NewStatement(entity.StatementSelect, dial)
	instance := &selectStmt{
		stmt: stmt,
		dial: dial,
	}

	stmt.LoadClauses([]entity.ClauseBuilder{
		&instance.clsSelect,
		&instance.clsFrom,
		&instance.clsJoin,
		&instance.clsWhere,
		&instance.clsGroupBy,
		&instance.clsHaving,
		&instance.clsOrderBy,
		&instance.clsLimit,
		&instance.clsOffset,
	})

	stmt.PreBuild(func(w *writer.WriterImpl) {
		if !instance.clsFrom.IsDeclared() && table != "" {
			instance.From(table)
		}

		if instance.clsJoin.IsDeclared() && table != "" {
			w.DefaultAlias(table)
		}
	})

	return instance
}

func (s *selectStmt) Test() (string, []any, error) {
	result := s.stmt.Build()
	return result.Query, result.Values, result.Err
}
