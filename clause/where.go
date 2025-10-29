package clause

import "strings"

type WhereClause struct {
	exprs []*Expr
}

func (w *WhereClause) Where(indetifier any) *Expr {
	expr := &Expr{}
	w.exprs = append(w.exprs, expr)
	return expr.Start(indetifier)
}

func (wc *WhereClause) Build(sb *strings.Builder, ph *Placeholder) {
	sb.WriteString("WHERE ")
	for i, expr := range wc.exprs {
		if i > 0 {
			sb.WriteByte(' ')
			sb.WriteString(string(and))
			sb.WriteByte(' ')
		}
		expr.build(sb, ph)
	}
}
