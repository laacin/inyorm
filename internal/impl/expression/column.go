package expression

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
)

type ColumnImpl struct{ dml.Column }

func (c *ColumnImpl) Count(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dml.ColExpr{
		Kind:  dml.ColAggrCount,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Sum(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dml.ColExpr{
		Kind:  dml.ColAggrSum,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Min(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dml.ColExpr{
		Kind:  dml.ColAggrMin,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Max(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dml.ColExpr{
		Kind:  dml.ColAggrMax,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Avg(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dml.ColExpr{
		Kind:  dml.ColAggrAvg,
		Value: dist,
	}
	return c
}

// --- Arith

func (c *ColumnImpl) Add(v any) api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind:  dml.ColArithAdd,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Sub(v any) api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind:  dml.ColArithSub,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Mul(v any) api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind:  dml.ColArithMul,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Div(v any) api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind:  dml.ColArithDiv,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Mod(v any) api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind:  dml.ColArithMod,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Wrap() api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind: dml.ColArithWrap,
	})
	return c
}

// --- Scalar

func (c *ColumnImpl) Lower() api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind: dml.ColScalarLower,
	})
	return c
}

func (c *ColumnImpl) Upper() api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind: dml.ColScalarUpper,
	})
	return c
}

func (c *ColumnImpl) Trim() api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind: dml.ColScalarTrim,
	})
	return c
}

func (c *ColumnImpl) Round() api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind: dml.ColScalarRound,
	})
	return c
}

func (c *ColumnImpl) Abs() api.Column {
	c.Exprs = append(c.Exprs, dml.ColExpr{
		Kind: dml.ColScalarAbs,
	})
	return c
}

// --- Alias

func (c *ColumnImpl) As(value string) api.Column {
	c.Alias = value
	return c
}

// --- Build

func (c *ColumnImpl) Build() dml.Value {
	return &c.Column
}
