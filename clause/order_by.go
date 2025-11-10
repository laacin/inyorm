package clause

import "github.com/laacin/inyorm/internal/core"

type OrderByClause struct {
	orders []*Order
}

func (o *OrderByClause) Name() core.ClauseType {
	return core.ClsTypOrderBy
}

func (o *OrderByClause) IsDeclared() bool { return o != nil }

func (o *OrderByClause) Build(w core.Writer) {
	w.Write("ORDER BY ")
	for i, ord := range o.orders {
		if i > 0 {
			w.Write(", ")
		}

		w.Value(ord.order, core.WriterOpts{ColType: core.ColTypAlias})
		if ord.descending {
			w.Write(" DESC")
		}
	}
}

// -- Methods

func (o *OrderByClause) OrderBy(value any) core.ClauseOrder {
	order := &Order{order: value}
	o.orders = append(o.orders, order)
	return order
}

// Depending clause

type Order struct {
	order      any
	descending bool
}

func (o *Order) Desc() { o.descending = true }
