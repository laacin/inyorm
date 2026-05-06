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
	col := entity.Column{Table: getTbl(expr.MainRef, ref), Name: name}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) All(ref ...string) Col {
	col := entity.Column{Table: getTbl(expr.MainRef, ref), From: &entity.Wildcard{}}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Param(value ...any) Param {
	store := len(value) > 0
	var v any
	if store {
		v = value[0]
	}

	impl := entity.Parameter{Store: store, Value: v}
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

	col := entity.Column{From: cs.Deref()}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Self, Col, Param, Cond, CondNext, CaseSwitch, CaseSearch, CaseNext,
]) Search(fn func(CaseSearch)) Col {
	cs := &CaseSearchImpl[CaseSwitch, CaseNext]{}
	fn(any(cs).(CaseSearch))

	col := entity.Column{From: cs.Deref()}
	impl := ColumnImpl[Col]{Column: col}
	return any(impl).(Col)

}

// --- Helpers
func getTbl(mainRef string, candidate []string) string {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return mainRef
}
