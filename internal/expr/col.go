package expr

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

// --- Entity

type Col struct {
	Name  string
	Ref   core.LazyVal[core.Reference]
	Alias string
	Value string
	from  Expr
	aggr  *ColAggr
	exprs []ColExpr
}

func NewCol(name string, ref core.LazyVal[core.Reference]) *Col {
	if ref == nil {
		return &Col{Name: name, Ref: func() core.Reference { return core.Reference{Enabled: false} }}
	}
	return &Col{Name: name, Ref: ref}
}
func NewColFrom(from Expr) *Col {
	return &Col{from: from}
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

func (*Col) Kind() Kind { return KindCol }

func (c *Col) Render(w core.InternalWriter, dial Renderer, mode core.WritingMode) {
	if c != nil && (c.aggr != nil || c.exprs != nil || c.from != nil) {
		nw := w.New()

		if c.from != nil {
			nw.Value(c.from, mode)
		} else {
			dial.WriteExprColBase(nw, c)
		}

		if c.exprs != nil {
			for _, e := range c.exprs {
				if scalar, ok := e.IsScalar(); ok {
					dial.WriteExprColScalar(nw, scalar)
					continue
				}

				if arith, ok := e.IsArith(); ok {
					dial.WriteExprColArith(nw, arith)
					continue
				}

				if e.IsWrap() {
					dial.WriteExprColWrap(nw)
					continue
				}
			}
			c.exprs = nil
		}

		if c.aggr != nil {
			dial.WriteExprColAggr(nw, c.aggr)
			c.aggr = nil
		}

		c.Value = nw.ToString()
	}

	switch mode {
	case core.WriteBase:
		dial.WriteExprColBase(w, c)

	case core.WriteExpr:
		dial.WriteExprColExpr(w, c)

	case core.WriteAlias:
		dial.WriteExprColAlias(w, c)

	case core.WriteDef:
		dial.WriteExprColDef(w, c)
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
