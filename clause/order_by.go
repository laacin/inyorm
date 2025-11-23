package clause

import "github.com/laacin/inyorm/internal/core"

type OrderBy[Next any] struct {
	declared bool
	orders   []*order
	current  *order
}

func (cls *OrderBy[Next]) Name() core.ClauseType { return core.ClsTypOrderBy }
func (cls *OrderBy[Next]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *OrderBy[Next]) Build(w core.Writer) {
	w.Write("ORDER BY")
	w.Char(' ')

	for i, ord := range cls.orders {
		if i > 0 {
			w.Write(", ")
		}
		w.Value(ord.order, cls.Name())
		if ord.descending {
			w.Write(" DESC")
		}
	}
}

// -- Methods

func (cls *OrderBy[Next]) OrderBy(value any) Next {
	cls.declared = true
	order := &order{order: value}
	cls.current = order
	cls.orders = append(cls.orders, order)
	return any(cls).(Next)
}

func (cls *OrderBy[Next]) Desc() { cls.current.descending = true }

// -- internal

type order struct {
	order      any
	descending bool
}
