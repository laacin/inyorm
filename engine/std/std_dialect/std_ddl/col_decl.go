package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

func (s *DdlSyntax) WriteTableDecl(w core.Writer, t *ddl.TableDecl) {
	w.Write("CREATE TABLE IF NOT EXISTS")
	w.Char(' ')
	w.Write(t.Name)

	lazyCons := []core.WriterFunc{}
	for _, c := range t.Cons {
		if cons, ok := c.IsForeignKey(); ok {
			lazyCons = append(lazyCons, func(w core.Writer) { s.WriteConsForeignKey(w, cons) })
			continue
		}
		if cons, ok := c.IsCheck(); ok {
			lazyCons = append(lazyCons, func(w core.Writer) { s.WriteConsCheck(w, cons) })
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

func (s *DdlSyntax) WriteColDecl(w core.Writer, c *ddl.ColDecl) {
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

	if c.Default != nil {
		w.Char(' ')
		s.WriteMetaDefault(w, c.Default)
	}
}

func (*DdlSyntax) WriteColText(w core.Writer) {
	w.Write("TEXT")
}
func (*DdlSyntax) WriteColInt(w core.Writer) {
	w.Write("INTEGER")
}
func (*DdlSyntax) WriteColFloat(w core.Writer) {
	w.Write("DOUBLE")
}
func (*DdlSyntax) WriteColBool(w core.Writer) {
	w.Write("BOOLEAN")
}

func (*DdlSyntax) WriteMetaPrimaryKey(w core.Writer) {
	w.Write("PRIMARY KEY")
}
func (*DdlSyntax) WriteMetaAutoIncrement(w core.Writer) {
	w.Write("AUTOINCREMENT")
}
func (*DdlSyntax) WriteMetaUnique(w core.Writer) {
	w.Write("UNIQUE")
}
func (*DdlSyntax) WriteMetaNotNull(w core.Writer) {
	w.Write("NOT NULL")
}
func (*DdlSyntax) WriteMetaDefault(w core.Writer, cons *ddl.ConsDefault) {
	w.Write("DEFAULT")
	w.Char(' ')
	w.Value(cons.Value, core.WriteBase)
}

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

func (*DdlSyntax) WriteConsIndex(w core.Writer, cons *ddl.ConsDecl[ddl.ConsIndex]) {
	w.Write("CREATE INDEX")
	w.Write(" ON ")

	w.Write(cons.Table)
	w.Char('(')
	w.Write(cons.Column)
	w.Char(')')
}

func (*DdlSyntax) WriteConsCheck(w core.Writer, cons *ddl.ConsDecl[ddl.ConsCheck]) {
	w.Write("CHECK")
	w.Char(' ')
	w.Value(cons.Value.Cond, core.WriteBase)
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
