package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

func (s *DllSyntax) WriteTableDecl(w core.Writer, t *ddl.TableDecl) {
	w.Write("CREATE TABLE IF NOT EXISTS")
	w.Char(' ')
	w.Write(t.Name)

	inlineCons := map[string]core.WriterFunc{}
	lazyCons := []core.WriterFunc{}
	outCons := []core.WriterFunc{}

	for _, c := range t.Cons {
		if cons, ok := c.IsForeignKey(); ok {
			lazyCons = append(lazyCons, func(w core.Writer) { s.WriteConsForeignKey(w, cons) })
			continue
		}

		if cons, ok := c.IsCheck(); ok {
			lazyCons = append(lazyCons, func(w core.Writer) { s.WriteConsCheck(w, cons) })
			continue
		}

		if cons, ok := c.IsIndex(); ok {
			outCons = append(outCons, func(w core.Writer) { s.WriteConsIndex(w, cons) })
			continue
		}

		if cons, ok := c.IsDefault(); ok {
			inlineCons[cons.Column] = func(w core.Writer) { s.WriteConsDefault(w, cons) }
			continue
		}
	}

	w.Char(' ')
	w.Char('(')
	for i, c := range t.Cols {
		if i > 0 {
			w.Write(", ")
		}
		s.WriteColDecl(w, &c)

		if cons, ok := inlineCons[c.Name]; ok {
			w.Char(' ')
			cons(w)
		}
	}

	for _, cons := range lazyCons {
		w.Write(", ")
		cons(w)
	}
	w.Write(");")

	if len(outCons) > 0 {
		w.Char('\n')

		for i, cons := range outCons {
			if i > 0 {
				w.Write(";\n")
			}
			cons(w)
		}
	}
}

func (s *DllSyntax) WriteColDecl(w core.Writer, c *ddl.ColDecl) {
	w.Write(c.Name)
	w.Char(' ')
	w.Write(mapColKind[c.Kind])

	if c.Meta.PrimaryKey {
		w.Char(' ')
		s.WriteMetaPrimaryKey(w)

		if c.Kind == ddl.ColKindInt && c.Meta.AutoIncrement {
			w.Char(' ')
			s.WriteMetaAutoIncrement(w)
		}

		return
	}

	if c.Meta.Unique {
		w.Char(' ')
		s.WriteMetaUnique(w)

		if c.Kind == ddl.ColKindInt && c.Meta.AutoIncrement {
			w.Char(' ')
			s.WriteMetaAutoIncrement(w)
		}
	}

	if c.Meta.NotNull {
		w.Char(' ')
		s.WriteMetaNotNull(w)
	}
}

func (*DllSyntax) WriteColText(w core.Writer) {
	w.Write("TEXT")
}
func (*DllSyntax) WriteColInt(w core.Writer) {
	w.Write("INTEGER")
}
func (*DllSyntax) WriteColFloat(w core.Writer) {
	w.Write("DOUBLE")
}
func (*DllSyntax) WriteColBool(w core.Writer) {
	w.Write("BOOLEAN")
}

func (*DllSyntax) WriteMetaPrimaryKey(w core.Writer) {
	w.Write("PRIMARY KEY")
}
func (*DllSyntax) WriteMetaAutoIncrement(w core.Writer) {
	w.Write("AUTOINCREMENT")
}
func (*DllSyntax) WriteMetaUnique(w core.Writer) {
	w.Write("UNIQUE")
}
func (*DllSyntax) WriteMetaNotNull(w core.Writer) {
	w.Write("NOT NULL")
}

func (*DllSyntax) WriteConsForeignKey(w core.Writer, cons *ddl.ConsDecl[ddl.ConsForeignKey]) {
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

func (*DllSyntax) WriteConsIndex(w core.Writer, cons *ddl.ConsDecl[ddl.ConsIndex]) {
	w.Write("CREATE INDEX")
	w.Write(" ON ")

	w.Write(cons.Table)
	w.Char('(')
	w.Write(cons.Column)
	w.Char(')')
}

func (*DllSyntax) WriteConsCheck(w core.Writer, cons *ddl.ConsDecl[ddl.ConsCheck]) {
	w.Write("CHECK")
	w.Char(' ')
	w.Value(cons.Value.Cond, core.WriteBase)
}

func (*DllSyntax) WriteConsDefault(w core.Writer, cons *ddl.ConsDecl[ddl.ConsDefault]) {
	w.Write("DEFAULT")
	w.Char(' ')
	w.Value(cons.Value.Value, core.WriteBase)
}

var mapOnAct = map[ddl.OnAction]string{
	ddl.OnActionCascade:  "CASCADE",
	ddl.OnActionSetNull:  "SET NULL",
	ddl.OnActionDefault:  "SET DEFAULT",
	ddl.OnActionRestrict: "RESTRICT",
	ddl.OnActionNoAction: "NO ACTION",
}

var mapColKind = map[ddl.ColKind]string{
	ddl.ColKindText:  "TEXT",
	ddl.ColKindInt:   "INTEGER",
	ddl.ColKindFloat: "DOUBLE",
	ddl.ColKindBool:  "BOOLEAN",
}
