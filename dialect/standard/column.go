package standard

import "github.com/laacin/inyorm/internal/entity"

func (dial *DialectStandard) WriteTable(w entity.Writer, tbl *entity.Table) {
	w.Write(tbl.Value)
	if ref, ok := w.GetRef(tbl.Value); ok {
		w.Char(' ')
		w.Char(ref)
	}
}

func (dial *DialectStandard) WriteColBase(w entity.Writer, col *entity.Column) {
	if col.Table != "" {
		if ref, ok := w.GetRef(col.Table); ok {
			w.Char(ref)
			w.Char('.')
		}
	}
	w.Write(col.Name)
}

func (dial *DialectStandard) WriteColExpr(w entity.Writer, col *entity.Column) {
	dial.BuildCol(w.New(), col)

	if col.Value == "" {
		dial.WriteColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (dial *DialectStandard) WriteColAlias(w entity.Writer, col *entity.Column) {
	dial.BuildCol(w.New(), col)

	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		dial.WriteColExpr(w, col)
		return
	}

	dial.WriteColBase(w, col)
}

func (dial *DialectStandard) WriteColDef(w entity.Writer, col *entity.Column) {
	dial.BuildCol(w.New(), col)

	if col.Value == "" {
		dial.WriteColBase(w, col)
		return
	}

	dial.WriteColExpr(w, col)
	if col.Alias != "" {
		w.Write(" AS ")
		w.Write(col.Alias)
	}
}

// --- Helpers
func (dial *DialectStandard) BuildFirst(w entity.Writer, col *entity.Column) {
	if col.From != nil {
		if col.Table != "" && col.From.Kind() == entity.ValueWildcard {
			if ref, ok := w.GetRef(col.Table); ok {
				w.Char(ref)
				w.Char('.')
			}
		}
		col.From.Write(w, dial, entity.WriteExpr)
		return
	}

	if col.Value != "" {
		w.Write(col.Value)
		return
	}

	dial.WriteColBase(w, col)
}

// FIX: illegible
func (dial *DialectStandard) BuildCol(w entity.Writer, col *entity.Column) {
	if (col == nil) || (col.Exprs == nil && col.Aggr == nil && col.From == nil) {
		return
	}

	dial.BuildFirst(w, col)
	if col.Exprs != nil {
		for _, expr := range col.Exprs {
			if scalar, ok := scalarMap[expr.Kind]; ok {
				wScalar(scalar)(w)
				continue
			}

			if arith, ok := arithMap[expr.Kind]; ok {
				wArith(arith, expr.Value)(w)
				continue
			}

			if expr.Kind == entity.ColArithWrap {
				wWrap()(w)
				continue
			}
		}
		col.Exprs = nil
	}

	if col.Aggr != nil {
		if aggr, ok := aggrMap[col.Aggr.Kind]; ok {
			wAggr(col.Aggr.Value, aggr)(w)
		}
		col.Aggr = nil
	}

	col.Value = w.Result()
}

// maps
var aggrMap = map[entity.ColKindExpr]string{
	entity.ColAggrCount: "COUNT",
	entity.ColAggrSum:   "SUM",
	entity.ColAggrMin:   "MIN",
	entity.ColAggrMax:   "MAX",
	entity.ColAggrAvg:   "AVG",
}

var scalarMap = map[entity.ColKindExpr]string{
	entity.ColScalarLower: "LOWER",
	entity.ColScalarUpper: "UPPER",
	entity.ColScalarTrim:  "TRIM",
	entity.ColScalarRound: "ROUND",
	entity.ColScalarAbs:   "ABS",
}

var arithMap = map[entity.ColKindExpr]byte{
	entity.ColArithAdd: '+',
	entity.ColArithSub: '-',
	entity.ColArithMul: '*',
	entity.ColArithDiv: '/',
	entity.ColArithMod: '%',
}

func wArith(arg byte, value any) entity.WriterFunc {
	return func(w entity.Writer) {
		w.Char(' ')
		w.Char(arg)
		w.Char(' ')
		w.Value(value, entity.WriteExpr)
	}
}

func wScalar(arg string) entity.WriterFunc {
	return func(w entity.Writer) {
		prev := w.Result()
		w.Reset()

		w.Write(arg)
		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
}

func wWrap() entity.WriterFunc {
	return func(w entity.Writer) {
		prev := w.Result()
		w.Reset()

		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
}

func wAggr(distinct any, aggr string) entity.WriterFunc {
	return func(w entity.Writer) {
		prev := w.Result()
		w.Reset()

		w.Write(aggr)
		w.Char('(')
		if dist, ok := distinct.(bool); ok && dist {
			w.Write("DISTINCT")
			w.Char(' ')
		}
		w.Write(prev)
		w.Char(')')
	}
}
