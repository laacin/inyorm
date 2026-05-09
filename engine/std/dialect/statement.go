package dialect

import "github.com/laacin/inyorm/internal/entity"

func (dial *DialectStd) SelectOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseSelect,
		entity.ClauseFrom,
		entity.ClauseJoin,
		entity.ClauseWhere,
		entity.ClauseGroupBy,
		entity.ClauseHaving,
		entity.ClauseOrderBy,
		entity.ClauseLimit,
		entity.ClauseOffset,
	}
}

func (dial *DialectStd) InsertOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseInsertInto,
	}
}

func (dial *DialectStd) UpdateOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseUpdate,
		entity.ClauseWhere,
	}
}

func (dial *DialectStd) DeleteOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseDelete,
		entity.ClauseFrom,
		entity.ClauseWhere,
	}
}
