package statement

import (
	"context"

	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/mapper"
	"github.com/laacin/inyorm/internal/query"
)

type Transaction struct {
	driver     core.Driver
	candidates []*TxStatement

	err error
}

func NewTransaction(driver core.Driver) *Transaction {
	if driver == nil {
		return &Transaction{err: errExecNoDriver}
	}

	return &Transaction{
		driver: driver,
		err:    nil,
	}
}

func (tx *Transaction) Push(qc *query.Compiler) api.SelfBinder {
	if tx.err != nil {
		return &TxStatement{}
	}

	result, err := qc.Compile()
	if err != nil {
		tx.err = err
		return &TxStatement{}
	}

	stmt := &TxStatement{query: result.QueryString}
	stmt.Binder = NewBinder[api.SelfBinder](result.Params, stmt)

	tx.candidates = append(tx.candidates, stmt)
	return stmt
}

func (tx *Transaction) Run(context ...context.Context) error {
	if tx.err != nil {
		return tx.err
	}

	ctx := core.OptionalCtx(context)
	t := tx.driver.BeginTx(ctx)

	for _, cand := range tx.candidates {
		vals, err := cand.params.Values()
		if err != nil {
			_ = t.Rollback()
			return err
		}

		if cand.bind == nil {
			if err := t.Exec(ctx, cand.query, vals...); err != nil {
				_ = t.Rollback()
				return err
			}
			continue
		}

		rows, err := t.Query(ctx, cand.query, vals...)
		if err != nil {
			_ = t.Rollback()
			return err
		}

		if err := mapper.New().Bind(rows, cand.bind); err != nil {
			_ = t.Rollback()
			return err
		}
	}

	_ = t.Commit()
	return nil
}

// helpers

type TxStatement struct {
	query string
	Binder[api.SelfBinder]
}
