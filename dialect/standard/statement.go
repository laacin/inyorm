package standard

import "github.com/laacin/inyorm/internal/entity"

func (dial *DialectStandard) SelectOrder() []entity.ClauseKind {
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

func (dial *DialectStandard) InsertOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseInsertInto,
	}
}

func (dial *DialectStandard) UpdateOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseUpdate,
		entity.ClauseWhere,
	}
}

func (dial *DialectStandard) DeleteOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseDelete,
		entity.ClauseWhere,
	}
}
