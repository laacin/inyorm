package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity
type Table struct{ Value string }

// --- Builder
type TableBuilder struct{ emb Table }

// start

func (t *TableBuilder) Start(table string) *TableBuilder {
	t.emb.Value = table
	return t
}

// --- Build
func (t *TableBuilder) Kind() ExprKind {
	return ExprTable
}

func (t *TableBuilder) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	w.SetRef(t.emb.Value)
	dial.WriteTable(w, &t.emb)
}
