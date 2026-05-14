package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type ColumnImpl struct{ expr.Column }

func (c *ColumnImpl) Count(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &expr.ColExpr{
		Kind:  expr.ColAggrCount,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Sum(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &expr.ColExpr{
		Kind:  expr.ColAggrSum,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Min(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &expr.ColExpr{
		Kind:  expr.ColAggrMin,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Max(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &expr.ColExpr{
		Kind:  expr.ColAggrMax,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Avg(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &expr.ColExpr{
		Kind:  expr.ColAggrAvg,
		Value: dist,
	}
	return c
}

// --- Arith

func (c *ColumnImpl) Add(v any) api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind:  expr.ColArithAdd,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Sub(v any) api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind:  expr.ColArithSub,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Mul(v any) api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind:  expr.ColArithMul,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Div(v any) api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind:  expr.ColArithDiv,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Mod(v any) api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind:  expr.ColArithMod,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Wrap() api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind: expr.ColArithWrap,
	})
	return c
}

// --- Scalar

func (c *ColumnImpl) Lower() api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind: expr.ColScalarLower,
	})
	return c
}

func (c *ColumnImpl) Upper() api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind: expr.ColScalarUpper,
	})
	return c
}

func (c *ColumnImpl) Trim() api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind: expr.ColScalarTrim,
	})
	return c
}

func (c *ColumnImpl) Round() api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind: expr.ColScalarRound,
	})
	return c
}

func (c *ColumnImpl) Abs() api.Column {
	c.Exprs = append(c.Exprs, expr.ColExpr{
		Kind: expr.ColScalarAbs,
	})
	return c
}

// --- Alias

func (c *ColumnImpl) As(value string) api.Column {
	c.Alias = value
	return c
}

// --- Build

func (c *ColumnImpl) Build() expr.Value {
	return &c.Column
}
