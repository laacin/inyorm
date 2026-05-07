package expression

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type ColumnImpl struct{ entity.Column }

func (c *ColumnImpl) Count(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrCount,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Sum(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrSum,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Min(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrMin,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Max(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrMax,
		Value: dist,
	}
	return c
}

func (c *ColumnImpl) Avg(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrAvg,
		Value: dist,
	}
	return c
}

// --- Arith

func (c *ColumnImpl) Add(v any) api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithAdd,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Sub(v any) api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithSub,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Mul(v any) api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithMul,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Div(v any) api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithDiv,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Mod(v any) api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithMod,
		Value: v,
	})
	return c
}

func (c *ColumnImpl) Wrap() api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColArithWrap,
	})
	return c
}

// --- Scalar

func (c *ColumnImpl) Lower() api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarLower,
	})
	return c
}

func (c *ColumnImpl) Upper() api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarUpper,
	})
	return c
}

func (c *ColumnImpl) Trim() api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarTrim,
	})
	return c
}

func (c *ColumnImpl) Round() api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarRound,
	})
	return c
}

func (c *ColumnImpl) Abs() api.Column {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarAbs,
	})
	return c
}

// --- Alias

func (c *ColumnImpl) As(value string) api.Column {
	c.Alias = value
	return c
}

// --- Build

func (c *ColumnImpl) Build() entity.Value {
	return &c.Column
}
