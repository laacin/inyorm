package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type TableImpl struct {
	emb expr.Table
}

func (t *TableImpl) Start(table string) api.Table {
	t.emb.Value = table
	return t
}

// --- Build
func (t *TableImpl) Kind() expr.ExprKind {
	return expr.ExprTable
}

func (t *TableImpl) Build(w core.InternalWriter, dial expr.ExprWriter, mode core.WritingMode) {
	w.SetRef(t.emb.Value)
	dial.WriteTable(w, &t.emb)
}
