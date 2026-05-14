package ddl

const (
	KeySep = ","
	ValSep = ":"
)

const (
	PrimaryKey     = "pk"
	PrimaryKeyLong = "primary_key"

	AutoIncrement     = "ai"
	AutoIncrementLong = "auto_increment"

	Unique     = "uq"
	UniqueLong = "unique"

	NotNull     = "nn"
	NotNullLong = "not_null"

	Index     = "idx"
	IndexLong = "index"

	Column     = "c"
	ColumnLong = "column"

	Foreign        = "fk"
	ForeignKeyLong = "foreign_key"

	OnDelete     = "od"
	OnDeleteLong = "on_delete"

	OnUpdate     = "ou"
	OnUpdateLong = "on_update"

	ActCascade  = "cascade"
	ActSetNull  = "setnull"
	ActDefault  = "default"
	ActRestrict = "restrict"
	ActNoAction = "noaction"
)

var (
	pkSet = map[string]struct{}{
		PrimaryKey:     {},
		PrimaryKeyLong: {},
	}

	aiSet = map[string]struct{}{
		AutoIncrement:     {},
		AutoIncrementLong: {},
	}

	uqSet = map[string]struct{}{
		Unique:     {},
		UniqueLong: {},
	}

	nnSet = map[string]struct{}{
		NotNull:     {},
		NotNullLong: {},
	}

	idxSet = map[string]struct{}{
		Index:     {},
		IndexLong: {},
	}

	cSet = map[string]struct{}{
		Column:     {},
		ColumnLong: {},
	}

	fkSet = map[string]struct{}{
		Foreign:        {},
		ForeignKeyLong: {},
	}

	odSet = map[string]struct{}{
		OnDelete:     {},
		OnDeleteLong: {},
	}

	ouSet = map[string]struct{}{
		OnUpdate:     {},
		OnUpdateLong: {},
	}
)
