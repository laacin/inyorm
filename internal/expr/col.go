package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type Col struct {
	Name  string
	Ref   string
	Alias string
	Value string
	from  ExprBuilder
	aggr  *ColAggr
	exprs []ColExpr
}

// start

func (c *Col) StartFrom(from ExprBuilder) *Col {
	c.from = from
	return c
}
func (c *Col) Start(name, table string) *Col {
	c.Name = name
	c.Ref = table
	return c
}

// --- PUB API

func (c *Col) Count(distinct ...bool) api.Col {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &ColAggr{
		Kind:     ColAggrCount,
		Distinct: dist,
	}
	return c
}

func (c *Col) Sum(distinct ...bool) api.Col {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &ColAggr{
		Kind:     ColAggrSum,
		Distinct: dist,
	}
	return c
}

func (c *Col) Min(distinct ...bool) api.Col {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &ColAggr{
		Kind:     ColAggrMin,
		Distinct: dist,
	}
	return c
}

func (c *Col) Max(distinct ...bool) api.Col {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &ColAggr{
		Kind:     ColAggrMax,
		Distinct: dist,
	}
	return c
}

func (c *Col) Avg(distinct ...bool) api.Col {
	dist := len(distinct) > 0 && distinct[0]
	c.aggr = &ColAggr{
		Kind:     ColAggrAvg,
		Distinct: dist,
	}
	return c
}

// --- Arith

func (c *Col) Add(v any) api.Col {
	c.exprs = append(c.exprs, newArith(ColArithAdd, v))
	return c
}

func (c *Col) Sub(v any) api.Col {
	c.exprs = append(c.exprs, newArith(ColArithSub, v))
	return c
}

func (c *Col) Mul(v any) api.Col {
	c.exprs = append(c.exprs, newArith(ColArithMul, v))
	return c
}

func (c *Col) Div(v any) api.Col {
	c.exprs = append(c.exprs, newArith(ColArithDiv, v))
	return c
}

func (c *Col) Mod(v any) api.Col {
	c.exprs = append(c.exprs, newArith(ColArithMod, v))
	return c
}

func (c *Col) Wrap() api.Col {
	c.exprs = append(c.exprs, newWrap())
	return c
}

// --- Scalar

func (c *Col) Lower() api.Col {
	c.exprs = append(c.exprs, newScalar(ColScalarLower))
	return c
}

func (c *Col) Upper() api.Col {
	c.exprs = append(c.exprs, newScalar(ColScalarUpper))
	return c
}

func (c *Col) Trim() api.Col {
	c.exprs = append(c.exprs, newScalar(ColScalarTrim))
	return c
}

func (c *Col) Round() api.Col {
	c.exprs = append(c.exprs, newScalar(ColScalarRound))
	return c
}

func (c *Col) Abs() api.Col {
	c.exprs = append(c.exprs, newScalar(ColScalarAbs))
	return c
}

// --- Alias

func (c *Col) As(value string) api.Col {
	c.Alias = value
	return c
}

// --- Build
func (c *Col) BaseName() string { return c.Name }

func (c *Col) Kind() ExprKind {
	return ExprCol
}

func (c *Col) Build(w core.InternalWriter, dial ExprWriter, mode core.WritingMode) {
	w.SetRef(c.Ref)

	if c != nil && (c.aggr != nil || c.exprs != nil || c.from != nil) {
		nw := w.New()

		if c.from != nil {
			nw.Value(c.from, mode)
		} else {
			dial.WriteColBase(nw, c)
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

		c.Value = nw.ToString()
	}

	switch mode {
	case core.WriteBase:
		dial.WriteColBase(w, c)

	case core.WriteExpr:
		dial.WriteColExpr(w, c)

	case core.WriteAlias:
		dial.WriteColAlias(w, c)

	case core.WriteDef:
		dial.WriteColDef(w, c)
	}
}

// helpers

func newArith(kind ColArithKind, value any) ColExpr {
	return ColExpr{
		Kind: ColKindArith,
		Value: ColArith{
			Kind:  kind,
			Value: value,
		},
	}
}

func newScalar(kind ColScalarKind) ColExpr {
	return ColExpr{
		Kind: ColKindScalar,
		Value: ColScalar{
			Kind: kind,
		},
	}
}

func newWrap() ColExpr {
	return ColExpr{Kind: ColKindWrap}
}
