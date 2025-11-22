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
	query := writer.NewQuery(table, cfg)
	stmt := &selectStmt{
		query: query,
		executor: &executor{
			Cfg:      cfg,
			Ctx:      ctx,
			Instance: db,
			Query:    query,
		},
	}

	clauses := []core.Clause{
		&stmt.clsSelect, &stmt.clsFrom, &stmt.clsJoin,
		&stmt.clsWhere, &stmt.clsGroupBy, &stmt.clsHaving,
		&stmt.clsOrderBy, &stmt.clsLimit, &stmt.clsOffset,
	}

	query.SetClauses(clauses, []core.ClauseType{
		core.ClsTypSelect,
		core.ClsTypFrom,
		core.ClsTypJoin,
		core.ClsTypWhere,
		core.ClsTypGroupBy,
		core.ClsTypHaving,
		core.ClsTypOrderBy,
		core.ClsTypLimit,
		core.ClsTypOffset,
	})

	return stmt
}
