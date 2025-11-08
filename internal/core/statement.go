package core

type SelectStatement interface {
	ClauseSelect
	ClauseFrom
	ClauseJoin
	ClauseWhere
	ClauseGroupBy
	ClauseOrderBy
	ClauseLimit
	ClauseOffset
	As(alias string)
}

type InsertStatement interface {
	ClauseInsert
}
