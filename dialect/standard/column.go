package standard

import "github.com/laacin/inyorm/intr/dialect"

func (dial *DialectStandard) Table(w dialect.Writer, tbl dialect.Table) {
	w.Write(tbl.Name)
	if ref, shouldBeUsed := w.GetTableRef(tbl.Name); shouldBeUsed {
		w.Char(' ')
		w.Char(ref)
	}
}

func (dial *DialectStandard) ColWriteBase(w dialect.Writer, col dialect.Column) {
	if col.Table != "" {
		if ref, shouldBeUsed := w.GetTableRef(col.Table); shouldBeUsed {
			w.Char(ref)
			w.Char('.')
		}
	}
	w.Write(col.Name)
}

func (dial *DialectStandard) ColWriteExpr(w dialect.Writer, col dialect.Column) {
	dial.BuildCol(w.New(), &col)

	if col.Value == "" {
		dial.ColWriteBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (dial *DialectStandard) ColWriteAlias(w dialect.Writer, col dialect.Column) {
	dial.BuildCol(w.New(), &col)

	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		dial.ColWriteExpr(w, col)
		return
	}

	dial.ColWriteBase(w, col)
}

func (dial *DialectStandard) ColWriteDef(w dialect.Writer, col dialect.Column) {
	dial.BuildCol(w.New(), &col)

	if col.Value == "" {
		dial.ColWriteBase(w, col)
		return
	}

	dial.ColWriteExpr(w, col)
	if col.Alias != "" {
		w.Write(" AS ")
		w.Write(col.Alias)
	}
}

// -- Custom columns
func (dial *DialectStandard) ColConcat(values []any) dialect.WriterFunc {
	return func(w dialect.Writer) {
		w.Write("CONCAT")
		w.Char('(')
		for i, val := range values {
			if i > 0 {
				w.Write(", ")
			}
			w.Value(val, dialect.WriteExpr)
		}
		w.Char(')')
	}
}

func (dial *DialectStandard) ColSwitch(cond any, cas dialect.CaseCond) dialect.WriterFunc {
	return func(w dialect.Writer) {
		w.Write("CASE")
		w.Char(' ')
		w.Value(cond, dialect.WriteExpr)

		for _, expr := range cas.Exprs {
			w.Write(" WHEN ")
			w.Value(expr.Identifier, dialect.WriteExpr)

			w.Write(" THEN ")
			w.Value(expr.Argument, dialect.WriteExpr)
			w.Char(' ')
		}

		if cas.Els != nil {
			w.Write("ELSE")
			w.Char(' ')
			w.Value(cas.Els, dialect.WriteExpr)
			w.Char(' ')
		}

		w.Write("END")
	}
}

func (dial *DialectStandard) ColSearch(cas dialect.CaseCond) dialect.WriterFunc {
	return func(w dialect.Writer) {
		w.Write("CASE WHEN")
		w.Char(' ')

		for _, arg := range cas.Exprs {
			dial.Cond(w, arg.Identifier.(dialect.Cond), dialect.WriteExpr) // NOTE: fragile
			w.Write(" THEN ")
			w.Value(arg.Argument, dialect.WriteExpr)
			w.Char(' ')
		}

		if cas.Els != nil {
			w.Write("ELSE")
			w.Char(' ')
			w.Value(cas.Els, dialect.WriteExpr)
			w.Char(' ')
		}

		w.Write("END")
	}
}

// --- Helpers
func (dial *DialectStandard) BuildCol(w dialect.Writer, col *dialect.Column) {
	if (col == nil) || (col.Exprs == nil && col.Aggr == nil && col.Complex == nil) {
		return
	}

	if col.Exprs != nil || col.Aggr != nil {
		func() {
			if col.Complex == nil && col.Value == "" {
				dial.ColWriteBase(w, *col)
				return
			}

			if col.Value != "" {
				w.Write(col.Value)
				return
			}

			if col.Complex != nil {
				col.Complex(w)
			}
		}()
	}

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

			if expr.Kind == dialect.ColArithWrap {
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
var aggrMap = map[dialect.ColKindExpr]string{
	dialect.ColAggrCount: "COUNT",
	dialect.ColAggrSum:   "SUM",
	dialect.ColAggrMin:   "MIN",
	dialect.ColAggrMax:   "MAX",
	dialect.ColAggrAvg:   "AVG",
}

var scalarMap = map[dialect.ColKindExpr]string{
	dialect.ColScalarLower: "LOWER",
	dialect.ColScalarUpper: "UPPER",
	dialect.ColScalarTrim:  "TRIM",
	dialect.ColScalarRound: "ROUND",
	dialect.ColScalarAbs:   "ABS",
}

var arithMap = map[dialect.ColKindExpr]byte{
	dialect.ColArithAdd: '+',
	dialect.ColArithSub: '-',
	dialect.ColArithMul: '*',
	dialect.ColArithDiv: '/',
	dialect.ColArithMod: '%',
}

func wArith(arg byte, value any) dialect.WriterFunc {
	return func(w dialect.Writer) {
		prev := w.Result()
		w.Reset()

		w.Write(prev)
		w.Char(' ')
		w.Char(arg)
		w.Char(' ')
		w.Value(value, dialect.WriteExpr)
	}
}

func wScalar(arg string) dialect.WriterFunc {
	return func(w dialect.Writer) {
		prev := w.Result()
		w.Reset()

		w.Write(arg)
		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
}

func wWrap() dialect.WriterFunc {
	return func(w dialect.Writer) {
		prev := w.Result()
		w.Reset()

		w.Char('(')
		w.Write(prev)
		w.Char(')')
	}
}

func wAggr(distinct any, aggr string) dialect.WriterFunc {
	return func(w dialect.Writer) {
		prev := w.Result()
		w.Reset()

		w.Write(aggr)
		w.Char('(')
		if _, ok := distinct.(bool); ok {
			w.Write("DISTINCT")
			w.Char(' ')
		}
		w.Write(prev)
		w.Char(')')
	}
}
