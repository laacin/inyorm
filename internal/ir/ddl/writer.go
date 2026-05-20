package ddl

import "github.com/laacin/inyorm/internal/core"

type Syntax interface {
	TableWriter
}

type TableWriter interface {
	WriteCreateTable(core.Writer, *TableDecl)
	WriteColDecl(core.Writer, *ColDecl)

	WriteColText(core.Writer)
	WriteColInt(core.Writer)
	WriteColFloat(core.Writer)
	WriteColBool(core.Writer)

	WriteMetaPrimaryKey(core.Writer)
	WriteMetaAutoIncrement(core.Writer)
	WriteMetaUnique(core.Writer)
	WriteMetaNotNull(core.Writer)
	WriteMetaDefault(core.Writer, any)

	WriteConsForeignKey(core.Writer, *ConsDecl[ConsForeignKey])
	WriteConsCheck(core.Writer, *ConsDecl[ConsCheck])
	WriteConsIndex(core.Writer, *ConsDecl[ConsIndex])
}

// --- Internal
type TableBuilder interface {
	Build(core.InternalWriter, TableWriter) error
}
