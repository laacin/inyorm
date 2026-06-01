package std_dialect

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query/ddl"
)

// ---- CREATE TABLE ----

func (s *Dialect) WriteQueryCreateTable(w core.Writer, t *ddl.QueryCreateTable) {
	w.Write("CREATE TABLE IF NOT EXISTS")
	w.Char(' ')
	w.Write(t.Name)

	w.Char(' ')
	w.Char('(')
	for i, c := range t.Cols {
		if i > 0 {
			w.Write(", ")
		}
		s.Self.WriteColDecl(w, c)
	}

	if t.ConsPk != nil {
		w.Write(", ")
		s.Self.WriteConsPrimaryKey(w, t.ConsPk)
	}

	for _, fk := range t.Fks {
		w.Write(", ")
		s.Self.WriteConsForeignKey(w, fk)
	}

	for _, ck := range t.Checks {
		w.Write(", ")
		s.Self.WriteConsCheck(w, ck)
	}
	w.Write(")")
}

// ---- WRITE COLUMN DECLARATION ----

func (s *Dialect) WriteColDecl(w core.Writer, c *ddl.ColDecl) {
	w.Write(c.Name)
	w.Char(' ')
	switch c.Kind {
	case ddl.ColKindString:
		s.Self.WriteColString(w)
	case ddl.ColKindInt:
		s.Self.WriteColInt(w)
	case ddl.ColKindFloat:
		s.Self.WriteColFloat(w)
	case ddl.ColKindBool:
		s.Self.WriteColBool(w)
	default:
		panic("unexpected col kind")
	}

	if c.Meta.PrimaryKey {
		w.Char(' ')
		s.Self.WriteMetaPrimaryKey(w)

		if c.Kind == ddl.ColKindInt && c.Meta.AutoIncrement {
			w.Char(' ')
			s.Self.WriteMetaAutoIncrement(w)
		}

		return
	}

	if c.Meta.Unique {
		w.Char(' ')
		s.Self.WriteMetaUnique(w)

		if c.Kind == ddl.ColKindInt && c.Meta.AutoIncrement {
			w.Char(' ')
			s.Self.WriteMetaAutoIncrement(w)
		}
	}

	if c.Meta.NotNull {
		w.Char(' ')
		s.Self.WriteMetaNotNull(w)
	}

	if c.Meta.Default != nil {
		w.Char(' ')
		s.Self.WriteMetaDefault(w, c.Meta.Default)
	}
}

// ---- COLUMN DECLARATION KINDS & META ----

func (*Dialect) WriteColString(w core.Writer) {
	w.Write("TEXT")
}
func (*Dialect) WriteColInt(w core.Writer) {
	w.Write("INTEGER")
}
func (*Dialect) WriteColFloat(w core.Writer) {
	w.Write("DOUBLE")
}
func (*Dialect) WriteColBool(w core.Writer) {
	w.Write("BOOLEAN")
}

func (*Dialect) WriteMetaPrimaryKey(w core.Writer) {
	w.Write("PRIMARY KEY")
}
func (*Dialect) WriteMetaAutoIncrement(w core.Writer) {
	w.Write("AUTOINCREMENT")
}
func (*Dialect) WriteMetaUnique(w core.Writer) {
	w.Write("UNIQUE")
}
func (*Dialect) WriteMetaNotNull(w core.Writer) {
	w.Write("NOT NULL")
}
func (*Dialect) WriteMetaDefault(w core.Writer, value any) {
	w.Write("DEFAULT")
	w.Char(' ')
	w.Value(value, core.WriteBase)
}

// ---- TABLE CONSTRAINTS ----

func (*Dialect) WriteConsPrimaryKey(w core.Writer, cons *ddl.PrimaryKey) {
	w.Write("PRIMARY KEY")
	w.Char(' ')
	w.Char('(')
	for i, col := range cons.Cols {
		if i > 0 {
			w.Write(", ")
		}
		w.Write(col)
	}
	w.Char(')')
}

func (*Dialect) WriteConsForeignKey(w core.Writer, cons *ddl.ForeignKey) {
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

func (*Dialect) WriteConsCheck(w core.Writer, cons *ddl.Check) {
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

// --- Helpers

var mapOnAct = map[ddl.OnAction]string{
	ddl.OnActionCascade:  "CASCADE",
	ddl.OnActionSetNull:  "SET NULL",
	ddl.OnActionDefault:  "SET DEFAULT",
	ddl.OnActionRestrict: "RESTRICT",
	ddl.OnActionNoAction: "NO ACTION",
}
