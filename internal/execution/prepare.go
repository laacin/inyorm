package execution

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/entity/constant"
)

type PrepareExec struct {
	ctx  context.Context
	stmt *sql.Stmt
}

func (e *PrepareExec) Run(args []any, binder ...any) error {
	qty := len(binder)

	if qty == 0 {
		return runPrep(e.ctx, e.stmt, args)
	}

	if qty == 1 {
		return scanPrep(e.ctx, e.stmt, constant.TAG, args, binder[0])
	}

	return scanPrep(e.ctx, e.stmt, constant.TAG, args, binder)
}

// -- Internal

func (e *PrepareExec) close() error {
	if err := e.stmt.Close(); err != nil {
		return errSQL(err)
	}
	return nil
}
