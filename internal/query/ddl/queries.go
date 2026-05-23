package ddl

import "github.com/laacin/inyorm/internal/core"

type CreateTable struct {
	Name   string
	Cols   []*ColDecl
	Fks    []*ForeignKey
	Checks []*Check
}

func (b *CreateTable) Build(w core.InternalWriter, dial TableWriter) {
	dial.WriteCreateTable(w, b)
}
