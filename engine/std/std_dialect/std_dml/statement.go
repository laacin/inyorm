package std_dml

import "github.com/laacin/inyorm/internal/query/dml"

func (*DmlSyntax) SelectOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseKindSelect,
		dml.ClauseKindFrom,
		dml.ClauseKindJoin,
		dml.ClauseKindWhere,
		dml.ClauseKindGroupBy,
		dml.ClauseKindHaving,
		dml.ClauseKindOrderBy,
		dml.ClauseKindLimit,
		dml.ClauseKindOffset,
	}
}

func (*DmlSyntax) InsertOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseKindInsert,
	}
}

func (*DmlSyntax) UpdateOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseKindUpdate,
		dml.ClauseKindWhere,
	}
}

func (*DmlSyntax) DeleteOrder() []dml.ClauseKind {
	return []dml.ClauseKind{
		dml.ClauseKindDelete,
		dml.ClauseKindFrom,
		dml.ClauseKindWhere,
	}
}
