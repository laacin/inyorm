package std_dml

import "github.com/laacin/inyorm/internal/entity/dml"

func (*DmlSyntax) SelectOrder() []dml.ClauseKind {
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

func (*DmlSyntax) InsertOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseInsertInto,
	}
}

func (*DmlSyntax) UpdateOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseUpdate,
		dml.ClauseWhere,
	}
}

func (*DmlSyntax) DeleteOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseDelete,
		dml.ClauseFrom,
		dml.ClauseWhere,
	}
}
