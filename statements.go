package inyorm

import (
	"context"

	"github.com/laacin/inyorm/internal"
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/impl/clause"
	"github.com/laacin/inyorm/internal/writer"
)

type selectStmt struct {
	stmt *internal.Statement
	dial entity.Dialect

	clause.SelectImpl
	clause.FromImpl
	clause.JoinImpl
	clause.WhereImpl
	clause.GroupByImpl
	clause.HavingImpl
	clause.OrderByImpl
	clause.LimitImpl
	clause.OffsetImpl

	// clsSelect
	// clsFrom
	// clsJoin
	// clsWhere
	// clsGroupBy
	// clsHaving
	// clsOrderBy
	// clsLimit
	// clsOffset
}

func newSelect(ctx context.Context, dial entity.Dialect, table string) *selectStmt {
	stmt := internal.NewStatement(entity.StatementSelect, dial)
	instance := &selectStmt{
		stmt: stmt,
		dial: dial,
	}

	stmt.LoadClauses([]entity.ClauseBuilder{
		&instance.SelectImpl,
		&instance.FromImpl,
		&instance.JoinImpl,
		&instance.WhereImpl,
		&instance.GroupByImpl,
		&instance.HavingImpl,
		&instance.OrderByImpl,
		&instance.LimitImpl,
		&instance.OffsetImpl,
	})

	stmt.PreBuild(func(w *writer.WriterImpl) {
		if !instance.FromImpl.IsDeclared() && table != "" {
			instance.From(&entity.Table{Value: table})
		}

		if instance.JoinImpl.IsDeclared() && table != "" {
			w.DefaultAlias(table)
		}
	})

	return instance
}

func (s *selectStmt) Test() (string, []any, error) {
	result := s.stmt.Build()
	return result.Query, result.Values, result.Err
}
