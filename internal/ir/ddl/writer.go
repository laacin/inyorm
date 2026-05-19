package ddl

import "github.com/laacin/inyorm/internal/core"

type TableWriter interface {
	WriteTableDecl(core.Writer, *TableDecl)
	WriteColDecl(core.Writer, *ColDecl)

	WriteColText(core.Writer)
	WriteColInt(core.Writer)
	WriteColFloat(core.Writer)
	WriteColBool(core.Writer)

	WriteMetaPrimaryKey(core.Writer)
	WriteMetaAutoIncrement(core.Writer)
	WriteMetaUnique(core.Writer)
	WriteMetaNotNull(core.Writer)

	WriteConsForeignKey(core.Writer, *ConsDecl[ConsForeignKey])
	WriteConsIndex(core.Writer, *ConsDecl[ConsIndex])
	WriteConsDefault(core.Writer, *ConsDecl[ConsDefault])
	WriteConsCheck(core.Writer, *ConsDecl[ConsCheck])
}

// --- Internal
type TableBuilder interface {
	Build(core.InternalWriter, TableWriter) error
}
