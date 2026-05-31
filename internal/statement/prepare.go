package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/mapper"
	"github.com/laacin/inyorm/internal/query"
)

type Prepared struct {
	stmt core.Prepared
	Binder[api.Prepared]
	snapshot core.ParamStore
}

func NewPrepared(ctx context.Context, driver core.Driver, qc *query.Compiler) (*Prepared, error) {
	result, err := qc.Compile()
	if err != nil {
		return nil, err
	}

	stmt, err := driver.Prepare(ctx, result.QueryString)
	if err != nil {
		return nil, err
	}
	self := &Prepared{stmt: stmt, snapshot: result.Params}
	self.Binder = NewBinder[api.Prepared](self.snapshot.Clone(), self)
	return self, nil
}

func (p *Prepared) Run(context ...context.Context) error {
	defer func() {
		p.params = p.snapshot.Clone()
		p.bind = nil
	}()

	if p.err != nil {
		return p.err
	}

	vals, err := p.params.Values()
	if err != nil {
		return err
	}

	ctx := core.OptionalCtx(context)
	if p.bind == nil {
		return p.stmt.Exec(ctx, vals...)
	}

	rows, err := p.stmt.Query(ctx, vals...)
	if err != nil {
		return err
	}

	return mapper.New().Bind(rows, p.bind)
}
