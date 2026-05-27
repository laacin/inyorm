package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity
type Table struct {
	Name string
	Ref  core.LazyVal[core.Reference]
}

// start

func (t *Table) Start(name string, ref core.LazyVal[core.Reference]) *Table {
	t.Name = name
	t.Ref = ref
	return t
}

// --- Build

func (t *Table) Kind() Kind { return KindTable }

func (t *Table) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteExprTable(w, t)
}
