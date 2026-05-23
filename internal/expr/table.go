package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity
type Table struct{ Value string }

// start

func (t *Table) Start(table string) *Table {
	t.Value = table
	return t
}

// --- Build
func (t *Table) Kind() ExprKind {
	return ExprTable
}

func (t *Table) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	w.SetRef(t.Value)
	dial.WriteTable(w, t)
}
