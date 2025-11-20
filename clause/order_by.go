package clause

import "github.com/laacin/inyorm/internal/core"

type OrderBy[Next, Ident any] struct {
	declared bool
	orders   []*order
	current  *order
}

func (cls *OrderBy[Next, Ident]) Name() core.ClauseType { return core.ClsTypOrderBy }
func (cls *OrderBy[Next, Ident]) IsDeclared() bool      { return cls != nil && cls.declared }
func (cls *OrderBy[Next, Ident]) Build(w core.Writer) {
	w.Write("ORDER BY")
	w.Char(' ')

	for i, ord := range cls.orders {
		if i > 0 {
			w.Write(", ")
		}
		w.Identifier(ord.order, cls.Name())
		if ord.descending {
			w.Write(" DESC")
		}
	}
}

// -- Methods

func (cls *OrderBy[Next, Ident]) OrderBy(value Ident) Next {
	cls.declared = true
	order := &order{order: value}
	cls.current = order
	cls.orders = append(cls.orders, order)
	return any(cls).(Next)
}

func (cls *OrderBy[Next, Ident]) Desc() { cls.current.descending = true }

// -- internal

type order struct {
	order      any
	descending bool
}
