package exec

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type Executor[Prep any] struct {
	Ctx      context.Context
	Cfg      *core.Config
	Instance *sql.DB
	Query    *writer.Query
}

func (e *Executor[Prep]) Raw() (string, []any) { return e.Query.Build() }

func (e *Executor[Prep]) Run() error {
	query, args := e.Query.Build()
	return run(e.Ctx, e.Instance, query, args)
}

func (e *Executor[Prep]) Find(binder any) error {
	query, args := e.Query.Build()
	return find(e.Ctx, e.Instance, e.Cfg.ColumnTag, query, args, binder)
}

func (e *Executor[Prep]) Prepare(fn func(exec Prep) error) error {
	query, _ := e.Query.Build()
	prep, err := e.Instance.PrepareContext(e.Ctx, query)
	if err != nil {
		return errSQL(err)
	}

	exec := &PrepareExec{
		ctx:  e.Ctx,
		cfg:  e.Cfg,
		stmt: prep,
	}

	if err := fn(any(exec).(Prep)); err != nil {
		return err
	}

	return exec.close()
}
