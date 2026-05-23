package ddl

import "github.com/laacin/inyorm/internal/core"

type TableWriter interface {
	WriteCreateTable(core.Writer, *CreateTable)
	WriteColDecl(core.Writer, *ColDecl)

	WriteColString(core.Writer)
	WriteColInt(core.Writer)
	WriteColFloat(core.Writer)
	WriteColBool(core.Writer)

	WriteMetaPrimaryKey(core.Writer)
	WriteMetaAutoIncrement(core.Writer)
	WriteMetaUnique(core.Writer)
	WriteMetaNotNull(core.Writer)
	WriteMetaDefault(core.Writer, any)

	WriteConsForeignKey(core.Writer, *ForeignKey)
	WriteConsCheck(core.Writer, *Check)
}
