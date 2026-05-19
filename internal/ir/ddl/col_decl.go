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
	Name    string
	Kind    ColKind
	Meta    ColMeta
	Default *ConsDefault
}

// --- Dependencies
type ColMeta struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	NotNull       bool
}
