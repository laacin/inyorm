package column

import "github.com/laacin/inyorm/intr/dialect"

type ColumnImpl[Self any] struct {
	dialect.Column
}

// --- Aggregation

func (c *ColumnImpl[Self]) Count(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dialect.ColExpr{
		Kind:  dialect.ColAggrCount,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Sum(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dialect.ColExpr{
		Kind:  dialect.ColAggrSum,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Min(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dialect.ColExpr{
		Kind:  dialect.ColAggrMin,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Max(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dialect.ColExpr{
		Kind:  dialect.ColAggrMax,
		Value: dist,
	}
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Avg(distinct ...bool) Self {
	dist := len(distinct) > 0 && distinct[0]
	c.Aggr = &dialect.ColExpr{
		Kind:  dialect.ColAggrAvg,
		Value: dist,
	}
	return any(c).(Self)
}

// --- Arith

func (c *ColumnImpl[Self]) Add(v any) Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind:  dialect.ColArithAdd,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Sub(v any) Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind:  dialect.ColArithSub,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Mul(v any) Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind:  dialect.ColArithMul,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Div(v any) Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind:  dialect.ColArithDiv,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Mod(v any) Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind:  dialect.ColArithMod,
		Value: v,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Wrap() Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind: dialect.ColArithWrap,
	})
	return any(c).(Self)
}

// --- Scalar

func (c *ColumnImpl[Self]) Lower() Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind: dialect.ColScalarLower,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Upper() Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind: dialect.ColScalarUpper,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Trim() Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind: dialect.ColScalarTrim,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Round() Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind: dialect.ColScalarRound,
	})
	return any(c).(Self)
}

func (c *ColumnImpl[Self]) Abs() Self {
	c.Exprs = append(c.Exprs, dialect.ColExpr{
		Kind: dialect.ColScalarAbs,
	})
	return any(c).(Self)
}

// --- Alias

func (c *ColumnImpl[Self]) As(value string) Self {
	c.Alias = value
	return any(c).(Self)
}
