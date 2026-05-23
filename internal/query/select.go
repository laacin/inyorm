package query

import (
	"errors"

	"github.com/laacin/inyorm/internal/impl/writer"
	"github.com/laacin/inyorm/internal/query/dml"
)

type SelectQuery struct {
	Ref  string
	Dial Dialect

	dml.ClauseSelect
	dml.ClauseFrom
	dml.ClauseJoin
	dml.ClauseWhere
	dml.ClauseGroupBy
	dml.ClauseHaving
	dml.ClauseOrderBy
	dml.ClauseLimit
	dml.ClauseOffset
}

// start

func (q *SelectQuery) Start(dial Dialect, ref string) *SelectQuery {
	q.Dial = dial
	q.Ref = ref
	return q
}

// --- Build

func (*SelectQuery) Kind() QueryKind {
	return QuerySelect
}

func (q *SelectQuery) Build() (*QueryResult, error) {
	// --- Guards

	if !q.ClauseSelect.IsDeclared() {
		return nil, errors.New("clause 'SELECT' must be declared")
	}
	if !q.ClauseFrom.IsDeclared() {
		return nil, errors.New("clause 'FROM' must be declared")
	}

	// --- Load clauses
	clauses := []dml.ClauseBuilder{
		&q.ClauseSelect,
		&q.ClauseFrom,
		&q.ClauseJoin,
		&q.ClauseWhere,
		&q.ClauseGroupBy,
		&q.ClauseHaving,
		&q.ClauseOrderBy,
		&q.ClauseLimit,
		&q.ClauseOffset,
	}

	clauseMap := make(map[dml.ClauseKind]dml.ClauseBuilder)
	for _, cls := range clauses {
		if cls.IsDeclared() {
			clauseMap[cls.Kind()] = cls
		}
	}

	// --- Declarate writers

	var (
		params  = &writer.ParamStore{}
		aliases *writer.AliasStore
	)

	// --- Set table references if Join exists

	if q.ClauseJoin.IsDeclared() {
		aliases = &writer.AliasStore{}
	}

	// --- Write the statement

	w := &writer.WriterImpl{
		Syntax:  q.Dial,
		Params:  params,
		Aliases: aliases,
	}

	w.SetRef(q.Ref)

	first := true
	for _, ord := range q.Dial.SelectOrder() {
		if cls, ok := clauseMap[ord]; ok {
			if !first {
				w.Char(' ')
			}
			first = false

			if err := cls.Build(w, q.Dial); err != nil {
				return nil, err
			}
		}
	}

	// --- Validate values

	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &QueryResult{
		Query:  w.ToString(),
		Values: params.Values(),
	}, nil
}
