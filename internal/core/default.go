package core

var DefaultColumnWriter = ColumnWriter{
	Select:  ColTypDef,
	Join:    ColTypBase,
	Where:   ColTypExpr,
	GroupBy: ColTypExpr,
	Having:  ColTypExpr,
	OrderBy: ColTypAlias,
}

const DefaultColumnTag = "col"
