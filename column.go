package inyorm

import (
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
)

type (
	Column struct {
		wrap core.Column
	}
	ColumnExpr struct {
		defaultTable string
		wrap         core.ColExpr
	}
	Builder = core.Builder
)

// ----- Constructors -----

func wrapColumn(wrap core.Column) *Column { return &Column{wrap} }

func newColExpr(defaultTable string) *ColumnExpr {
	c := &column.ColExpr{}
	return &ColumnExpr{defaultTable: defaultTable, wrap: c}
}

// ----- Column -----

func (c *Column) Def() Builder   { return c.wrap.Def }
func (c *Column) Expr() Builder  { return c.wrap.Expr }
func (c *Column) Alias() Builder { return c.wrap.Alias }
func (c *Column) Base() Builder  { return c.wrap.Base }

func (c *Column) Count(distinct ...bool) *Column {
	dist := len(distinct) > 0 && distinct[0]
	c.wrap.Count(dist)
	return c
}
func (c *Column) Sum(distinct ...bool) *Column {
	dist := len(distinct) > 0 && distinct[0]
	c.wrap.Sum(dist)
	return c
}
func (c *Column) Min(distinct ...bool) *Column {
	dist := len(distinct) > 0 && distinct[0]
	c.wrap.Min(dist)
	return c
}
func (c *Column) Max(distinct ...bool) *Column {
	dist := len(distinct) > 0 && distinct[0]
	c.wrap.Max(dist)
	return c
}
func (c *Column) Avg(distinct ...bool) *Column {
	dist := len(distinct) > 0 && distinct[0]
	c.wrap.Avg(dist)
	return c
}

func (c *Column) Add(v Value) *Column { c.wrap.Add(vOne(v)); return c }
func (c *Column) Sub(v Value) *Column { c.wrap.Sub(vOne(v)); return c }
func (c *Column) Mul(v Value) *Column { c.wrap.Mul(vOne(v)); return c }
func (c *Column) Div(v Value) *Column { c.wrap.Div(vOne(v)); return c }
func (c *Column) Mod(v Value) *Column { c.wrap.Mod(vOne(v)); return c }
func (c *Column) Wrap() *Column       { c.wrap.Wrap(); return c }

func (c *Column) Lower() *Column { c.wrap.Lower(); return c }
func (c *Column) Upper() *Column { c.wrap.Upper(); return c }
func (c *Column) Trim() *Column  { c.wrap.Trim(); return c }
func (c *Column) Round() *Column { c.wrap.Round(); return c }
func (c *Column) Abs() *Column   { c.wrap.Abs(); return c }

func (c *Column) As(alias string) *Column { c.wrap.As(alias); return c }

// ----- Column expression -----

func (c *ColumnExpr) Col(name string, table ...string) *Column {
	tbl := c.defaultTable
	if len(table) > 0 {
		tbl = table[0]
	}
	col := c.wrap.Col(name, tbl)
	return wrapColumn(col)
}

func (c *ColumnExpr) All() *Column {
	return wrapColumn(c.wrap.All())
}

func (c *ColumnExpr) Concat(values ...Value) *Column {
	return wrapColumn(c.wrap.Concat(vMany(values)))
}

func (c *ColumnExpr) Cond(identifier Value) *Condition {
	cond := c.wrap.Condition(vOne(identifier))
	return wrapCondition(cond)
}

func (c *ColumnExpr) Switch(on Value, fn func(cs CaseSwitch)) *Column {
	cs := newCase[Value]()
	fn(cs)
	col := c.wrap.Switch(on, cs.wrap)
	return wrapColumn(col)
}

func (c *ColumnExpr) Search(fn func(cs CaseSearch)) *Column {
	cs := newCase[*CondNext]()
	fn(cs)
	col := c.wrap.Search(cs.wrap)
	return wrapColumn(col)
}
