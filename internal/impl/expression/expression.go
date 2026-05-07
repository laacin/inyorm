package expression

import "github.com/laacin/inyorm/internal/entity"

type ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext any,
] struct {
	MainRef string
	Dialect entity.ValueWriter
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) Col(name string, ref ...string) Col {
	col := entity.Column{Table: getLast(expr.MainRef, ref), Name: name}
	impl := &ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) All(ref ...string) Col {
	col := entity.Column{Table: getLast(expr.MainRef, ref), From: &entity.Wildcard{}}
	impl := &ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) Param(value ...any) Param {
	impl := &entity.Parameter{Store: len(value) > 0, Value: getLast(nil, value)}
	return any(impl).(Param)
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) Cond(ident any) Cond {
	cond := &ConditionImpl[Cond, CondNext]{}
	return cond.Start(ident)
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) Concat(values ...any) Col {
	col := entity.Column{From: &entity.Concat{Values: values}}
	impl := &ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) Switch(cond any, fn func(Case)) Col {
	cs := &CaseSwitchImpl[Case, CaseNext]{}
	fn(any(cs.Start(cond)).(Case))

	col := entity.Column{From: cs.Build()}
	impl := &ColumnImpl[Col]{Column: col}
	return any(impl).(Col)
}

func (expr *ExpressionImpl[
	Col, Param, Cond, CondNext, Case, CaseNext,
]) Search(fn func(Case)) Col {
	cs := &CaseSearchImpl[Case, CaseNext]{}
	fn(any(cs).(Case))

	col := entity.Column{From: cs.Build()}
	impl := &ColumnImpl[Col]{Column: col}
	return any(impl).(Col)

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
