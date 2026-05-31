package expr

import "github.com/laacin/inyorm/internal/core"

// --- Entity
type Table struct {
	Name string
	Ref  core.LazyVal[core.Reference]
}

func NewTable(name string, ref core.LazyVal[core.Reference]) *Table {
	return &Table{Name: name, Ref: ref}
}

// --- Build

func (t *Table) Kind() Kind { return KindTable }

func (t *Table) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	dial.WriteExprTable(w, t)
}
