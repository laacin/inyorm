package std_ddl

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

// --- MAIN WRITER

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

	if c.Meta.Default != nil {
		w.Char(' ')
		s.WriteMetaDefault(w, c.Meta.Default)
	}
}

// --- KINDS

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

// --- META

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
func (*DdlSyntax) WriteMetaDefault(w core.Writer, value any) {
	w.Write("DEFAULT")
	w.Char(' ')
	w.Value(value, core.WriteBase)
}
