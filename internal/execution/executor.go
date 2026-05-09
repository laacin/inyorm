package execution

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/constant"
)

type Executor struct {
	Ctx       context.Context
	Driver    entity.Driver
	Statement entity.StatementBuilder
}

func (e *Executor) Run(scanner ...any) error {
	qty := len(scanner)
	result, err := e.Statement.Build()
	if err != nil {
		return err
	}

	if qty == 0 {
		return run(e.Ctx, e.Driver, result.Query, result.Values)
	}

	if qty == 1 {
		return scan(e.Ctx, e.Driver, constant.TAG, result.Query, result.Values, scanner[0])
	}

	return scan(e.Ctx, e.Driver, constant.TAG, result.Query, result.Values, scanner)
}

// TODO:

// func (e *Executor) Prepare(fn func(exec *PrepareExec) error) error {
// 	panic("TODO")
//
// 	result, err := e.Statement.Build()
// 	if err != nil {
// 		return err
// 	}
//
// 	prep, err := e.Instance.PrepareContext(e.Ctx, result.Query)
// 	if err != nil {
// 		return errSQL(err)
// 	}
//
// 	exec := &PrepareExec{
// 		ctx:  e.Ctx,
// 		stmt: prep,
// 	}
//
// 	if err := fn(exec); err != nil {
// 		return err
// 	}
//
// 	return exec.close()
// }
