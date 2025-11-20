package column

import "github.com/laacin/inyorm/internal/core"

func (c *Column[Self]) Expr() core.Builder {
	return func(w core.Writer) {
		c.Builder.build(w.Split(), c)

		if c.value == "" {
			w.Column(c.Table, c.BaseName)
			return
		}
		w.Write(c.value)
	}
}

func (c *Column[Self]) Alias() core.Builder {
	return func(w core.Writer) {
		c.Builder.build(w.Split(), c)

		if c.alias != "" {
			w.Write(c.alias)
			return
		}

		if c.value == "" {
			w.Column(c.Table, c.BaseName)
			return
		}
		w.Write(c.value)
	}
}

func (c *Column[Self]) Def() core.Builder {
	return func(w core.Writer) {
		c.Builder.build(w.Split(), c)

		if c.value == "" {
			w.Column(c.Table, c.BaseName)
			return
		}

		w.Write(c.value)
		if c.alias != "" {
			w.Write(" AS ")
			w.Write(c.alias)
		}
	}
}

func (c *Column[Self]) Base() core.Builder {
	return func(w core.Writer) {
		w.Column(c.Table, c.BaseName)
	}
}
