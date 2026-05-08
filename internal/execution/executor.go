package execution

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/constant"
)

type Executor struct {
	Ctx       context.Context
	Instance  *sql.DB
	Statement entity.StatementBuilder
}

func (e *Executor) Run(binder ...any) error {
	qty := len(binder)
	result, err := e.Statement.Build()
	if err != nil {
		return err
	}

	if qty == 0 {
		return run(e.Ctx, e.Instance, result.Query, result.Values)
	}

	if qty == 1 {
		return scan(e.Ctx, e.Instance, constant.TAG, result.Query, result.Values, binder[0])
	}

	return scan(e.Ctx, e.Instance, constant.TAG, result.Query, result.Values, binder)
}

func (e *Executor) Prepare(fn func(exec *PrepareExec) error) error {
	panic("TODO")

	result, err := e.Statement.Build()
	if err != nil {
		return err
	}

	prep, err := e.Instance.PrepareContext(e.Ctx, result.Query)
	if err != nil {
		return errSQL(err)
	}

	exec := &PrepareExec{
		ctx:  e.Ctx,
		stmt: prep,
	}

	if err := fn(exec); err != nil {
		return err
	}

	return exec.close()
}
