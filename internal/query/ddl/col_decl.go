package ddl

import "github.com/laacin/inyorm/internal/api"

type ColDecl struct {
	Name string
	Kind ColKind
	Meta ColMeta
}

func NewColDecl(name string, kind ColKind) *ColDecl {
	return &ColDecl{
		Name: name,
		Kind: kind,
		Meta: ColMeta{NotNull: true},
	}
}
func (b *ColDecl) PrimaryKey() api.ColDeclNext {
	b.Meta.PrimaryKey = true
	return b
}
func (b *ColDecl) AutoIncrement() api.ColDeclNext {
	b.Meta.AutoIncrement = true
	return b
}
func (b *ColDecl) Unique() api.ColDeclNext {
	b.Meta.Unique = true
	return b
}
func (b *ColDecl) Nullable() api.ColDeclNext {
	b.Meta.NotNull = false
	return b
}
func (b *ColDecl) Default(v any) api.ColDeclNext {
	b.Meta.Default = v
	return b
}

// --- Build

func (b *ColDecl) Build() error {
	return nil
}

// --- Tools

type ColKind int

const (
	ColKindString ColKind = iota
	ColKindInt
	ColKindFloat
	ColKindBool
)

type ColMeta struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	NotNull       bool
	Default       any
}
