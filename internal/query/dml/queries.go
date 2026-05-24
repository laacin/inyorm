package dml

import "github.com/laacin/inyorm/internal/core"

// ---- SELECT QUERY -----

type QuerySelect struct {
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

func (q *QuerySelect) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQuerySelect(w, q)
	return nil
}

// ----- INSERT QUERY ------

type QueryInsert struct {
	ClauseInsertInto
}

func (q *QueryInsert) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryInsert(w, q)
	return nil
}

// ----- UPDATE QUERY ------

type QueryUpdate struct {
	ClauseUpdate
	ClauseWhere
}

func (q *QueryUpdate) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryUpdate(w, q)
	return nil
}

// ----- DELETE QUERY ------

type QueryDelete struct {
	ClauseDelete
	ClauseFrom
	ClauseWhere
}

func (q *QueryDelete) Render(w core.InternalWriter, dial QueryWriter) error {
	dial.WriteQueryDelete(w, q)
	return nil
}
