package inyorm

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/query"
	"github.com/laacin/inyorm/internal/statement"
)

type DB struct {
	eng *Engine
	Queries[Statement]
}

func New(eng *Engine, err error) (*DB, error) {
	if err != nil {
		return nil, err
	}

	qb := newQueryBuilder(eng.Dialect, func(qc *query.Compiler) Statement {
		return statement.New(eng.Driver, qc)
	})

	return &DB{
		eng:     eng,
		Queries: qb,
	}, nil
}

// --- Transaction
func (db *DB) StartTx() Transaction {
	tx := &txBuilder{
		Transaction: statement.NewTransaction(db.eng.Driver),
	}

	qb := newQueryBuilder(db.eng.Dialect, func(qc *query.Compiler) api.SelfBinder {
		return tx.Push(qc)
	})

	tx.Queries = qb
	return tx
}

// --- Connection

func (db *DB) Close() error {
	return db.eng.Driver.Close()
}

// helpers

type txBuilder struct {
	*statement.Transaction
	Queries[api.SelfBinder]
}
