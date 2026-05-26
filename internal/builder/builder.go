package builder

import (
	"github.com/laacin/inyorm/internal/builder/aliases"
	"github.com/laacin/inyorm/internal/builder/mapper"
	"github.com/laacin/inyorm/internal/builder/params"
	"github.com/laacin/inyorm/internal/builder/writer"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

type Builder struct {
	w *writer.Writer
	p *params.ParamStore
	a *aliases.AliasStore
	m *mapper.Mapper
}

func New(dial expr.ExprWriter) *Builder {
	return &Builder{
		w: writer.New(dial),
		p: params.New(),
		a: aliases.New(),
		m: mapper.New(),
	}
}

// --- Tools

func (b *Builder) Writer() core.InternalWriter { return b.w }
func (b *Builder) Params() core.ParamStore     { return b.p }
func (b *Builder) Aliases() core.AliasStore    { return b.a }
func (b *Builder) Mapper() core.Mapper         { return b.m }

// --- Rules

// must be called one time per build
func (b *Builder) SetMainRef(name string) {
	b.a.SetMain(name)
}
