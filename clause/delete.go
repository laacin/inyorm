package clause

import "github.com/laacin/inyorm/internal/core"

type Delete struct {
	declared bool
}

func (cls *Delete) Name() core.ClauseType                 { return core.ClsTypDelete }
func (cls *Delete) IsDeclared() bool                      { return cls != nil && cls.declared }
func (cls *Delete) Build(w core.Writer, cfg *core.Config) { w.Write("DELETE") }

// -- Methods

func (cls *Delete) Delete() { cls.declared = true }
