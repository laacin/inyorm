package inyorm

import (
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type Engine struct {
	cfg *core.Config
	db  *sql.DB
}

func New(dialect string, db *sql.DB, opts *Options) *Engine {
	cfg := resolveOpts(dialect, &opts)
	return &Engine{cfg: cfg, db: db}
}

// ---- Statements

func (e *Engine) NewSelect(table string) (SelectStmt, ColumnBuilder) {
	query := writer.NewQuery(table, e.cfg)

	stmt := &SelectStatement{query: query}
	colBldr := ColumnBuilder(&colBuilder{Table: table})

	return stmt, colBldr
}
