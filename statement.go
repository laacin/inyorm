package inyorm

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type selectStmt struct {
	query *writer.Query
	clsSelect
	clsFrom
	clsJoin
	clsWhere
	clsGroupBy
	clsHaving
	clsOrderBy
	clsLimit
	clsOffset
	*executor
}

func newSelect(ctx context.Context, cfg *core.Config, db *sql.DB, table string) SelectStmt {
	query := &writer.Query{Config: cfg}
	stmt := &selectStmt{
		query: query,
		executor: &executor{
			Cfg:      cfg,
			Ctx:      ctx,
			Instance: db,
			Query:    query,
		},
	}

	query.SetClauses([]core.Clause{
		&stmt.clsSelect, &stmt.clsFrom, &stmt.clsJoin,
		&stmt.clsWhere, &stmt.clsGroupBy, &stmt.clsHaving,
		&stmt.clsOrderBy, &stmt.clsLimit, &stmt.clsOffset,
	})

	query.PreBuild(func(cfg *core.Config) (useAliases bool) {
		if !stmt.clsFrom.IsDeclared() && table != "" {
			stmt.From(table)
		}

		if cfg.Limit > 0 && !stmt.clsLimit.IsDeclared() {
			stmt.clsLimit.Limit(cfg.Limit)
		}

		return stmt.clsJoin.IsDeclared()
	})

	return stmt
}
