package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

func (*DdlSyntax) WriteConsForeignKey(w core.Writer, cons *ddl.ConsDecl[ddl.ConsForeignKey]) {
	w.Write("FOREIGN KEY")
	w.Char(' ')

	w.Char('(')
	w.Write(cons.Column)
	w.Char(')')

	w.Write(" REFERENCES ")

	w.Write(cons.Value.ToTable)
	w.Char('(')
	w.Write(cons.Value.ToColumn)
	w.Char(')')

	if cons.Value.OnUpdate != ddl.OnActionUnset {
		w.Write(" ON UPDATE ")
		w.Write(mapOnAct[cons.Value.OnUpdate])
	}

	if cons.Value.OnDelete != ddl.OnActionUnset {
		w.Write(" ON DELETE ")
		w.Write(mapOnAct[cons.Value.OnDelete])
	}
}

func (*DdlSyntax) WriteConsCheck(w core.Writer, cons *ddl.ConsDecl[ddl.ConsCheck]) {
	w.Write("CHECK")
	w.Char(' ')
	w.Value(cons.Value.Cond, core.WriteBase)
}

func (*DdlSyntax) WriteConsIndex(w core.Writer, cons *ddl.ConsDecl[ddl.ConsIndex]) {
	w.Write("CREATE INDEX")
	w.Write(" ON ")

	w.Write(cons.Table)
	w.Char('(')
	w.Write(cons.Column)
	w.Char(')')
}

