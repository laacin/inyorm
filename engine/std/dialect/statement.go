package dialect

import "github.com/laacin/inyorm/internal/entity"

func (dial *StdDialect) SelectOrder() []entity.ClauseKind {
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

func (dial *StdDialect) InsertOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseInsertInto,
	}
}

func (dial *StdDialect) UpdateOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseUpdate,
		entity.ClauseWhere,
	}
}

func (dial *StdDialect) DeleteOrder() []entity.ClauseKind {
	return []entity.ClauseKind{
		entity.ClauseDelete,
		entity.ClauseFrom,
		entity.ClauseWhere,
	}
}
