package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

type ColDecl struct {
	Name string
	Kind ColKind
	Meta ColMeta
}

// start

func (b *ColDecl) Start(name string, kind ColKind) *ColDecl {
	b.Name = name
	b.Kind = kind
	b.Meta.NotNull = true
	return b
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

func (b *ColDecl) Build(w core.InternalWriter, dial TableWriter) {
	dial.WriteColDecl(w, b)
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
