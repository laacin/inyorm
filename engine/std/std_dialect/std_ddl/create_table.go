package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

func (s *DdlSyntax) WriteCreateTable(w core.Writer, t *ddl.TableDecl) {
	w.Write("CREATE TABLE IF NOT EXISTS")
	w.Char(' ')
	w.Write(t.Name)

	w.Char(' ')
	w.Char('(')
	for i, c := range t.Cols {
		if i > 0 {
			w.Write(", ")
		}
		s.WriteColDecl(w, &c)
	}

	for _, c := range t.Cons {
		w.Write(", ")

		if cons, ok := c.IsForeignKey(); ok {
			s.WriteConsForeignKey(w, cons)
			continue
		}

		if cons, ok := c.IsCheck(); ok {
			s.WriteConsCheck(w, cons)
			continue
		}
	}
	w.Write(")")
}
