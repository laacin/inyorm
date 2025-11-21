package inyorm

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type selectStmt struct {
	query *writer.Query
	clsSelect
	clsFrom
	clsJoin
	clsWhere
	clsGroupBy
	clsHaving
	clsOrderBy
	clsLimit
	clsOffset
}

func (stmt *selectStmt) Builder() *writer.Query {
	clauses := []core.Clause{
		&stmt.clsSelect, &stmt.clsFrom, &stmt.clsJoin,
		&stmt.clsWhere, &stmt.clsGroupBy, &stmt.clsHaving,
		&stmt.clsOrderBy, &stmt.clsLimit, &stmt.clsOffset,
	}

	stmt.query.SetClauses(clauses, []core.ClauseType{
		core.ClsTypSelect,
		core.ClsTypFrom,
		core.ClsTypJoin,
		core.ClsTypWhere,
		core.ClsTypGroupBy,
		core.ClsTypHaving,
		core.ClsTypOrderBy,
		core.ClsTypLimit,
		core.ClsTypOffset,
	})

	return stmt.query
}
