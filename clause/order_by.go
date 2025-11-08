package clause

import "github.com/laacin/inyorm/internal/core"

type OrderByClause struct {
	orders []*Order
}

func (o *OrderByClause) Name() string {
	return core.ClsOrderBy
}

func (o *OrderByClause) Build() core.Builder {
	return func(w core.Writer) {
		w.Write("ORDER BY ")
		for i, ord := range o.orders {
			if i > 0 {
				w.Write(", ")
			}

			w.Value(ord.order, nil)
			if ord.descending {
				w.Write(" DESC")
			}
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
