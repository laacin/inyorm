package exec

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
)

type PrepareExec struct {
	ctx  context.Context
	cfg  *core.Config
	stmt *sql.Stmt
}

func (e *PrepareExec) Run(values []any) error {
	return runPrep(e.ctx, e.stmt, values)
}

func (e *PrepareExec) Find(values []any, binder any) error {
	return findPrep(e.ctx, e.stmt, e.cfg.ColumnTag, values, binder)
}

// -- Internal

func (e *PrepareExec) close() error {
	if err := e.stmt.Close(); err != nil {
		return errSQL(err)
	}
	return nil
}
