package expression

import "github.com/laacin/inyorm/internal/entity"

type ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext any,
] struct {
	MainRef string
	Dialect entity.ValueWriter
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Col(name string, ref ...string) Col {
	col := entity.Column{Table: getLast(expr.MainRef, ref), Name: name}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) All(ref ...string) Col {
	col := entity.Column{Table: getLast(expr.MainRef, ref), From: &entity.Wildcard{}}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Param(value ...any) Param {
	impl := entity.Parameter{Store: len(value) > 0, Value: getLast(nil, value)}
	return any(impl).(Param)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Cond(ident any) CondNext {
	cond := &ConditionImpl[Cond, CondNext]{}
	cond.Start(ident)
	return any(cond).(CondNext)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Concat(values ...any) Col {
	col := entity.Column{From: &entity.Concat{Values: values}}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Switch(cond any, fn func(CaseSwitch)) Col {
	cs := &CaseSwitchImpl[CaseSwitch, CaseNext]{}
	fn(any(cs.Start(cond)).(CaseSwitch))

	col := entity.Column{From: cs.Build()}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Search(fn func(CaseSearch)) Col {
	cs := &CaseSearchImpl[CaseSwitch, CaseNext]{}
	fn(any(cs).(CaseSearch))

	col := entity.Column{From: cs.Build()}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
