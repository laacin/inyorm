package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/ddl"
)

func (*DdlSyntax) WriteConsForeignKey(w core.Writer, cons *ddl.ForeignKey) {
	w.Write("FOREIGN KEY")
	w.Char(' ')

	w.Char('(')
	w.Write(cons.Col)
	w.Char(')')

	w.Write(" REFERENCES ")

	w.Write(cons.ToTable)
	w.Char('(')
	w.Write(cons.ToCol)
	w.Char(')')

	if cons.OnUpdate != ddl.OnActionUnset {
		w.Write(" ON UPDATE ")
		w.Write(mapOnAct[cons.OnUpdate])
	}

	if cons.OnDelete != ddl.OnActionUnset {
		w.Write(" ON DELETE ")
		w.Write(mapOnAct[cons.OnDelete])
	}
}

func (*DdlSyntax) WriteConsCheck(w core.Writer, cons *ddl.Check) {
	w.Write("CHECK")
	w.Char(' ')
	w.Value(cons.Cond, core.WriteBase)
}

// func (*DdlSyntax) WriteConsIndex(w core.Writer, cons *ddl.ConsDecl[ddl.ConsIndex]) {
// 	w.Write("CREATE INDEX")
// 	w.Write(" ON ")
//
// 	w.Write(cons.Table)
// 	w.Char('(')
// 	w.Write(cons.Column)
// 	w.Char(')')
// }
//
