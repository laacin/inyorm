package exprimpl

import "github.com/laacin/inyorm/internal/api"

type ExprBuilderImpl struct{ DefaultRef string }

func (e *ExprBuilderImpl) Table(name string) api.Table {
	tbl := &TableImpl{}
	return tbl.Start(name)
}

func (e *ExprBuilderImpl) Col(name string, ref ...string) api.Column {
	col := &ColumnImpl{}
	return col.Start(name, getLast(e.DefaultRef, ref))
}

func (e *ExprBuilderImpl) All(ref ...string) api.Column {
	col := &ColumnImpl{}
	return col.Start("*", getLast(e.DefaultRef, ref))
}

func (e *ExprBuilderImpl) Param(value ...any) api.Parameter {
	param := &ParameterImpl{}
	return param.Start(len(value) > 0, getLast(nil, value))
}

func (e *ExprBuilderImpl) Cond(ident any) api.Condition {
	cond := &ConditionImpl{}
	return cond.Start(ident)
}

func (e *ExprBuilderImpl) Concat(values ...any) api.Column {
	col := &ColumnImpl{}
	return col.StartFrom((&ConcatImpl{}).Start(values))
}

func (e *ExprBuilderImpl) Switch(cond any, fn func(api.Case)) api.Column {
	cs := &CaseSwitchImpl{}
	fn(cs)

	col := ColumnImpl{}
	return col.StartFrom(cs)
}

func (e *ExprBuilderImpl) Search(fn func(api.Case)) api.Column {
	cs := &CaseSearchImpl{}
	fn(cs)

	col := &ColumnImpl{}
	return col.StartFrom(cs)

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
