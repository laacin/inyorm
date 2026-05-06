package expression

import "github.com/laacin/inyorm/internal/entity"

type ColumnImpl[Self any] struct{ entity.Column }

func (c *ColumnImpl[Self]) Count(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrCount,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Sum(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrSum,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Min(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrMin,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Max(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrMax,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Avg(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &entity.ColExpr{
		Kind:  entity.ColAggrAvg,
		Value: dist,
	}
	return any(c).(Self)
}

// --- Arith

func (c *ColumnImpl[Self]) Add(v any) Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithAdd,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Sub(v any) Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithSub,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Mul(v any) Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithMul,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Div(v any) Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithDiv,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Mod(v any) Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind:  entity.ColArithMod,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Wrap() Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColArithWrap,
	})
	return any(c).(Self)
}

// --- Scalar

func (c *ColumnImpl[Self]) Lower() Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarLower,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Upper() Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarUpper,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Trim() Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarTrim,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Round() Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarRound,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Abs() Self {
	c.Exprs = append(c.Exprs, entity.ColExpr{
		Kind: entity.ColScalarAbs,
	})
	return any(c).(Self)
}

// --- Alias

func (c *ColumnImpl[Self]) As(value string) Self {
	c.Alias = value
	return any(c).(Self)
}

// --- Build

func (c *ColumnImpl[Self]) Build() entity.Value {
	return &c.Column
}
