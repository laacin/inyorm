package inyorm

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
)

type Engine struct {
	cfg core.Config
	db  *sql.DB
}

func New(dialect string, db *sql.DB, opts *Options) *Engine {
	cfg := resolveOpts(dialect, opts)
	return &Engine{cfg: *cfg, db: db}
}

// ---- Statements

func (e *Engine) NewSelect(ctx context.Context, table string) (SelectStmt, ColumnBuilder) {
	stmt := newSelect(ctx, &e.cfg, e.db, table)
	colBldr := ColumnBuilder(&colBuilder{Table: table})
	return stmt, colBldr
}

func (e *Engine) NewInsert(ctx context.Context, table string) (InsertStmt, ColumnBuilder) {
	stmt := newInsert(ctx, &e.cfg, e.db, table)
	colBldr := ColumnBuilder(&colBuilder{Table: table})
	return stmt, colBldr
}

func (e *Engine) NewUpdate(ctx context.Context, table string) (UpdateStmt, ColumnBuilder) {
	stmt := newUpdate(ctx, &e.cfg, e.db, table)
	colBldr := ColumnBuilder(&colBuilder{Table: table})
	return stmt, colBldr
}

func (e *Engine) NewDelete(ctx context.Context, table string) (DeleteStmt, ColumnBuilder) {
	stmt := newDelete(ctx, &e.cfg, e.db, table)
	colBldr := ColumnBuilder(&colBuilder{Table: table})
	return stmt, colBldr
}
