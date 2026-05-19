package table

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/ddl"
)

type ColDeclImpl struct {
	emb      ddl.ColDecl
	nullable bool
}

func (c *ColDeclImpl) Start(name string, kind ddl.ColKind) api.ColDecl {
	c.emb.Name = name
	c.emb.Kind = kind
	c.emb.Meta.NotNull = true
	return c
}

func (c *ColDeclImpl) PrimaryKey() api.ColDecl {
	c.emb.Meta.PrimaryKey = true
	return c
}
func (c *ColDeclImpl) Unique() api.ColDecl {
	c.emb.Meta.Unique = true
	return c
}
func (c *ColDeclImpl) AutoIncrement() api.ColDecl {
	c.emb.Meta.AutoIncrement = true
	return c
}
func (c *ColDeclImpl) Nullable() api.ColDecl {
	c.nullable = true
	c.emb.Meta.NotNull = false
	return c
}
func (c *ColDeclImpl) Default(value any) api.ColDecl {
	c.emb.Default = &ddl.ConsDefault{Value: value}
	return c
}

// --- Build

func (c *ColDeclImpl) Build(w core.InternalWriter, dial ddl.TableWriter) error {
	c.emb.Meta.NotNull = !c.nullable

	dial.WriteColDecl(w, &c.emb)
	return nil
}
