package clause

import "github.com/laacin/inyorm/internal/core"

type Limit struct {
	declared bool
	limit    int
}

func (cls *Limit) Name() string     { return "LIMIT" }
func (cls *Limit) IsDeclared() bool { return cls != nil && cls.declared }
func (cls *Limit) Build(w core.Writer, cfg *core.Config) error {
	lim := cls.limit
	if cfg.MaxLimit > 0 {
		lim = max(lim, cfg.MaxLimit)
	}

	w.Write("LIMIT")
	w.Char(' ')
	w.Value(lim, core.ColTypUnset)
	return nil
}

// -- Methods

func (cls *Limit) Limit(value int) {
	if value > 0 {
		cls.declared = true
		cls.limit = value
	}
}
