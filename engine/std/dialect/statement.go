package dialect

import "github.com/laacin/inyorm/internal/entity/dml"

func (dial *StdDialect) SelectOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseSelect,
		dml.ClauseFrom,
		dml.ClauseJoin,
		dml.ClauseWhere,
		dml.ClauseGroupBy,
		dml.ClauseHaving,
		dml.ClauseOrderBy,
		dml.ClauseLimit,
		dml.ClauseOffset,
	}
}

func (dial *StdDialect) InsertOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseInsertInto,
	}
}

func (dial *StdDialect) UpdateOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseUpdate,
		dml.ClauseWhere,
	}
}

func (dial *StdDialect) DeleteOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseDelete,
		dml.ClauseFrom,
		dml.ClauseWhere,
	}
}
