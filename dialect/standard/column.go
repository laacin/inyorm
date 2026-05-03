package standard

import (
	"errors"

	"github.com/laacin/inyorm/intr/dialect"
	"github.com/laacin/inyorm/intr/writer"
)

func (dial *DialectStandard) Table(w dialect.Writer, tbl dialect.Table, def bool) {
	w.Write(tbl.Name)
	if def {
		if ref, shouldBeUsed := w.GetTableRef(tbl.Name); shouldBeUsed {
			w.Char(' ')
			w.Char(ref)
		}
	}
}

func (dial *DialectStandard) ColBase(w dialect.Writer, col dialect.Column) {
	if col.Table != "" {
		if ref, shouldBeUsed := w.GetTableRef(col.Table); shouldBeUsed {
			w.Char(ref)
			w.Char('.')
		}
	}
	w.Write(col.Name)
}

func (dial *DialectStandard) ColExpr(w dialect.Writer, col dialect.Column) {
	if col.Value == "" {
		dial.ColBase(w, col)
		return
	}

	w.Write(col.Value)
}

func (dial *DialectStandard) ColAlias(w dialect.Writer, col dialect.Column) {
	if col.Alias != "" {
		w.Write(col.Alias)
		return
	}

	if col.Value != "" {
		dial.ColExpr(w, col)
		return
	}

	dial.ColBase(w, col)
}

func (dial *DialectStandard) ColDef(w dialect.Writer, col dialect.Column) {
	if col.Value == "" {
		dial.ColBase(w, col)
		return
	}

	dial.ColExpr(w, col)
	if col.Alias != "" {
		w.Write(" AS ")
		w.Write(col.Alias)
	}
}

// -- Essentials
func (dial *DialectStandard) BuildColExpr(exprs []dialect.ColExpr) (string, error) {
	w := writer.Writer{}

	var (
		current  string
		lastAggr string
	)

	for _, expr := range exprs {
		if aggr, exists := aggrMap[expr.Kind]; exists {
			lastAggr = aggr
			continue
		}

		if scalar, exists := scalarMap[expr.Kind]; exists {
			w.Write(scalar)
			w.Char('(')
			w.Write(expr.Current)
			w.Char(')')

			current = w.Result()
			w.Reset()
			continue
		}

		if op, exists := arithMap[expr.Kind]; exists {
			if expr.Value == nil {
				return "", errors.New("arithmetical operation with nil value")
			}

			w.Write(expr.Current)
			w.Char(' ')
			w.Char(op)
			w.Value(expr.Value, dialect.ClauseNameNone)

			current = w.Result()
			w.Reset()
			continue
		}

		if expr.Kind == dialect.ColArithWrap {
			w.Char('(')
			w.Write(expr.Current)
			w.Char(')')

			current = w.Result()
			w.Reset()
			continue
		}

		return "", errors.New("unknown column expression kind")
	}

	if lastAggr != "" {
		w.Write(lastAggr)
		w.Char('(')
		w.Write(current)
		w.Char(')')

		return w.Result(), nil
	}
	return current, nil
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
