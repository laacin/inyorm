package writer

import "github.com/laacin/inyorm/internal/core"

var SelectOrder = []core.ClauseType{
	core.ClsTypSelect,
	core.ClsTypFrom,
	core.ClsTypJoin,
	core.ClsTypWhere,
	core.ClsTypGroupBy,
	core.ClsTypHaving,
	core.ClsTypOrderBy,
	core.ClsTypLimit,
	core.ClsTypOffset,
}

var InsertOrder = []core.ClauseType{
	core.ClsTypInsertInto,
}
