package exprimpl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type ColumnImpl struct {
	emb   expr.Column
	from  expr.ExprBuilder
	aggr  *expr.ColAggr
	exprs []expr.ColExpr
}

// new
func (c *ColumnImpl) StartFrom(from expr.ExprBuilder) api.Column {
	c.from = from
	return c
}
func (c *ColumnImpl) Start(name, ref string) api.Column {
	c.emb.Name = name
	c.emb.Ref = ref
	return c
}

func (c *ColumnImpl) Count(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &expr.ColAggr{
		Kind:     expr.ColAggrCount,
		Distinct: dist,
	}
	return c
}

func (c *ColumnImpl) Sum(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &expr.ColAggr{
		Kind:     expr.ColAggrSum,
		Distinct: dist,
	}
	return c
}

func (c *ColumnImpl) Min(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &expr.ColAggr{
		Kind:     expr.ColAggrMin,
		Distinct: dist,
	}
	return c
}

func (c *ColumnImpl) Max(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &expr.ColAggr{
		Kind:     expr.ColAggrMax,
		Distinct: dist,
	}
	return c
}

func (c *ColumnImpl) Avg(distinct ...bool) api.Column {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &expr.ColAggr{
		Kind:     expr.ColAggrAvg,
		Distinct: dist,
	}
	return c
}

// --- Arith

func (c *ColumnImpl) Add(v any) api.Column {
	c.exprs = append(c.exprs, newArith(expr.ColArithAdd, v))
	return c
}

func (c *ColumnImpl) Sub(v any) api.Column {
	c.exprs = append(c.exprs, newArith(expr.ColArithSub, v))
	return c
}

func (c *ColumnImpl) Mul(v any) api.Column {
	c.exprs = append(c.exprs, newArith(expr.ColArithMul, v))
	return c
}

func (c *ColumnImpl) Div(v any) api.Column {
	c.exprs = append(c.exprs, newArith(expr.ColArithDiv, v))
	return c
}

func (c *ColumnImpl) Mod(v any) api.Column {
	c.exprs = append(c.exprs, newArith(expr.ColArithMod, v))
	return c
}

func (c *ColumnImpl) Wrap() api.Column {
	c.exprs = append(c.exprs, newWrap())
	return c
}

// --- Scalar

func (c *ColumnImpl) Lower() api.Column {
	c.exprs = append(c.exprs, newScalar(expr.ColScalarLower))
	return c
}

func (c *ColumnImpl) Upper() api.Column {
	c.exprs = append(c.exprs, newScalar(expr.ColScalarUpper))
	return c
}

func (c *ColumnImpl) Trim() api.Column {
	c.exprs = append(c.exprs, newScalar(expr.ColScalarTrim))
	return c
}

func (c *ColumnImpl) Round() api.Column {
	c.exprs = append(c.exprs, newScalar(expr.ColScalarRound))
	return c
}

func (c *ColumnImpl) Abs() api.Column {
	c.exprs = append(c.exprs, newScalar(expr.ColScalarAbs))
	return c
}

// --- Alias

func (c *ColumnImpl) As(value string) api.Column {
	c.emb.Alias = value
	return c
}

// --- Build
func (c *ColumnImpl) BaseName() string { return c.emb.Name }

func (c *ColumnImpl) Kind() expr.ExprKind {
	return expr.ExprColumn
}

func (c *ColumnImpl) Build(w core.InternalWriter, dial expr.ExprWriter, mode core.WritingMode) {
	w.SetRef(c.emb.Ref)

	if c != nil && (c.aggr != nil || c.exprs != nil || c.from != nil) {
		nw := w.New()

		if c.from != nil {
			nw.Value(c.from, mode)
		} else {
			dial.WriteColBase(nw, &c.emb)
		}

		if c.exprs != nil {
			for _, e := range c.exprs {
				if scalar, ok := e.IsScalar(); ok {
					dial.WriteColScalar(nw, scalar)
					continue
				}

				if arith, ok := e.IsArith(); ok {
					dial.WriteColArith(nw, arith)
					continue
				}

				if e.IsWrap() {
					dial.WriteColWrap(nw)
					continue
				}
			}
			c.exprs = nil
		}

		if c.aggr != nil {
			dial.WriteColAggr(nw, c.aggr)
			c.aggr = nil
		}

		c.emb.Value = nw.ToString()
	}

	switch mode {
	case core.WriteBase:
		dial.WriteColBase(w, &c.emb)

	case core.WriteExpr:
		dial.WriteColExpr(w, &c.emb)

	case core.WriteAlias:
		dial.WriteColAlias(w, &c.emb)

	case core.WriteDef:
		dial.WriteColDef(w, &c.emb)
	}
}

// helpers
func newArith(kind expr.ColArithKind, value any) expr.ColExpr {
	return expr.ColExpr{
		Kind: expr.ColKindArith,
		Value: expr.ColArith{
			Kind:  kind,
			Value: value,
		},
	}
}

func newScalar(kind expr.ColScalarKind) expr.ColExpr {
	return expr.ColExpr{
		Kind: expr.ColKindScalar,
		Value: expr.ColScalar{
			Kind: kind,
		},
	}
}

func newWrap() expr.ColExpr {
	return expr.ColExpr{Kind: expr.ColKindWrap}
}
