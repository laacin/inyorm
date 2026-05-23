package dml

import "github.com/laacin/inyorm/internal/core"

// ---- SELECT QUERY -----

type SelectQuery struct {
	ClauseSelect
	ClauseFrom
	ClauseJoin
	ClauseWhere
	ClauseGroupBy
	ClauseHaving
	ClauseOrderBy
	ClauseLimit
	ClauseOffset
}

func (q *SelectQuery) Build(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQuerySelect(w, q)
	return nil
}

// ----- INSERT QUERY ------

type InsertQuery struct {
	ClauseInsert
}

func (q *InsertQuery) Build(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryInsert(w, q)
	return nil
}

// ----- UPDATE QUERY ------

type UpdateQuery struct {
	ClauseUpdate
	ClauseWhere
}

func (q *UpdateQuery) Build(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryUpdate(w, q)
	return nil
}

// ----- DELETE QUERY ------

type DeleteQuery struct {
	ClauseDelete
	ClauseFrom
	ClauseWhere
}

func (q *DeleteQuery) Build(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryDelete(w, q)
	return nil
}
