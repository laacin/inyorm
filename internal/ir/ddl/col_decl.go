package ddl

// --- Column type
type ColKind int

const (
	ColKindText ColKind = iota
	ColKindInt
	ColKindFloat
	ColKindBool
)

type ColDecl struct {
	Name string
	Kind ColKind
	Meta ColMeta
}

// --- Dependencies
type ColMeta struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	NotNull       bool
	Default       any
}
