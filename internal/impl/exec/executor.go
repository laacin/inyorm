package exec

import (
	"context"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/driver"
)

type Executor struct {
	Ctx       context.Context
	Driver    driver.Driver
	Statement interface{ Build() (string, []any, error) }
}

func (e *Executor) Run(scanner ...any) error {
	qty := len(scanner)
	query, values, err := e.Statement.Build()
	if err != nil {
		return err
	}

	if qty == 0 {
		return run(e.Ctx, e.Driver, query, values)
	}

	if qty == 1 {
		return scan(e.Ctx, e.Driver, core.TAG, query, values, scanner[0])
	}

	return scan(e.Ctx, e.Driver, core.TAG, query, values, scanner)
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
