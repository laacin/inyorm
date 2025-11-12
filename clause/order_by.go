package clause

import "github.com/laacin/inyorm/internal/core"

type OrderBy struct {
	orders  []*order
	current *order
}

func (o *OrderBy) Name() core.ClauseType { return core.ClsTypOrderBy }
func (o *OrderBy) IsDeclared() bool      { return o != nil }
func (o *OrderBy) Build(w core.Writer) {
	w.Write("ORDER BY")
	w.Char(' ')

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

func (o *OrderBy) OrderBy(value any) {
	order := &order{order: value}
	o.current = order
	o.orders = append(o.orders, order)
}

func (o *OrderBy) Desc() { o.current.descending = true }

// -- internal

type order struct {
	order      any
	descending bool
}
