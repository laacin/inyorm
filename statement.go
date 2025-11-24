package inyorm

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type selectStmt struct {
	q *writer.Query
	*executor

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

type insertStmt struct {
	q *writer.Query
	*executor

	clsInsert
}

type updateStmt struct {
	q *writer.Query
	*executor

	clsUpdate
	clsWhere
}

type deleteStmt struct {
	q *writer.Query
	*executor

	clsDelete
	clsFrom
	clsWhere
}

func newSelect(ctx context.Context, cfg *core.Config, db *sql.DB, table string) SelectStmt {
	q := &writer.Query{Config: cfg}
	exec := &executor{Ctx: ctx, Cfg: cfg, Instance: db, Query: q}
	stmt := &selectStmt{q: q, executor: exec}

	q.SetClauses([]core.Clause{
		&stmt.clsSelect,
		&stmt.clsFrom,
		&stmt.clsJoin,
		&stmt.clsWhere,
		&stmt.clsGroupBy,
		&stmt.clsHaving,
		&stmt.clsOrderBy,
		&stmt.clsLimit,
		&stmt.clsOffset,
	})

	q.PreBuild(func(cfg *core.Config) (useAliases bool) {
		if !stmt.clsFrom.IsDeclared() && table != "" {
			stmt.From(table)
		}

		if cfg.Limit > 0 && !stmt.clsLimit.IsDeclared() {
			stmt.Limit(cfg.Limit)
		}

		return stmt.clsJoin.IsDeclared()
	})

	return stmt
}

func newInsert(ctx context.Context, cfg *core.Config, db *sql.DB, table string) InsertStmt {
	q := &writer.Query{Config: cfg}
	exec := &executor{Ctx: ctx, Cfg: cfg, Instance: db, Query: q}
	stmt := &insertStmt{q: q, executor: exec}

	q.SetClauses([]core.Clause{
		&stmt.clsInsert,
	})

	q.PreBuild(func(cfg *core.Config) (useAliases bool) {
		stmt.Into(table)
		return false
	})

	return stmt
}

func newUpdate(ctx context.Context, cfg *core.Config, db *sql.DB, table string) UpdateStmt {
	q := &writer.Query{Config: cfg}
	exec := &executor{Ctx: ctx, Cfg: cfg, Instance: db, Query: q}
	stmt := &updateStmt{q: q, executor: exec}

	q.SetClauses([]core.Clause{
		&stmt.clsUpdate,
		&stmt.clsWhere,
	})

	q.PreBuild(func(cfg *core.Config) (useAliases bool) {
		stmt.clsUpdate.To(table)
		return false
	})

	return stmt
}

func newDelete(ctx context.Context, cfg *core.Config, db *sql.DB, table string) DeleteStmt {
	q := &writer.Query{Config: cfg}
	exec := &executor{Ctx: ctx, Cfg: cfg, Instance: db, Query: q}
	stmt := &deleteStmt{q: q, executor: exec}

	q.SetClauses([]core.Clause{
		&stmt.clsDelete,
		&stmt.clsFrom,
		&stmt.clsWhere,
	})

	q.PreBuild(func(cfg *core.Config) (useAliases bool) {
		stmt.clsDelete.Delete()

		if !stmt.clsFrom.IsDeclared() && table != "" {
			stmt.From(table)
		}

		return false
	})

	return stmt
}
