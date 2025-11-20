package core

var DefaultAutoPlaceholder = AutoPlaceholder{
	Insert: true,
	Update: true,
	Where:  true,
	Having: false,
	Join:   false,
}

var DefaultColumnWriter = ColumnWriter{
	Select:  ColTypDef,
	Join:    ColTypBase,
	Where:   ColTypExpr,
	GroupBy: ColTypExpr,
	Having:  ColTypExpr,
	OrderBy: ColTypAlias,
}

const DefaultColumnTag = "col"
