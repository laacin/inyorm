package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/ddl"
)

func (s *DdlSyntax) WriteCreateTable(w core.Writer, t *ddl.CreateTable) {
	w.Write("CREATE TABLE IF NOT EXISTS")
	w.Char(' ')
	w.Write(t.Name)

	w.Char(' ')
	w.Char('(')
	for i, c := range t.Cols {
		if i > 0 {
			w.Write(", ")
		}
		s.WriteColDecl(w, c)
	}

	for _, fk := range t.Fks {
		w.Write(", ")
		s.WriteConsForeignKey(w, fk)
	}

	for _, ck := range t.Checks {
		w.Write(", ")
		s.WriteConsCheck(w, ck)
	}
	w.Write(")")
}
