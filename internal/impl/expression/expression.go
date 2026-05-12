package expression

import (
	"github.com/laacin/inyorm/internal/entity/api"
	"github.com/laacin/inyorm/internal/entity/dml"
)

type ExprBuilderImpl struct{ DefaultRef string }

func (e *ExprBuilderImpl) Table(name string) api.Table {
	return &dml.Table{Value: name}
}

func (e *ExprBuilderImpl) Col(name string, ref ...string) api.Column {
	col := dml.Column{Ref: getLast(e.DefaultRef, ref), Name: name}
	return &ColumnImpl{Column: col}
}

func (e *ExprBuilderImpl) All(ref ...string) api.Column {
	col := dml.Column{Ref: getLast(e.DefaultRef, ref), From: &dml.Wildcard{}}
	return &ColumnImpl{Column: col}
}

func (e *ExprBuilderImpl) Param(value ...any) api.Parameter {
	return &dml.Parameter{Store: len(value) > 0, Value: getLast(nil, value)}
}

func (e *ExprBuilderImpl) Cond(ident any) api.Condition {
	cond := &ConditionImpl{}
	return cond.Start(ident)
}

func (e *ExprBuilderImpl) Concat(values ...any) api.Column {
	col := dml.Column{From: &dml.Concat{Values: values}}
	return &ColumnImpl{Column: col}
}

func (e *ExprBuilderImpl) Switch(cond any, fn func(api.Case)) api.Column {
	cs := &CaseSwitchImpl{}
	fn(cs)

	col := dml.Column{From: cs.Build()}
	return &ColumnImpl{Column: col}
}

func (e *ExprBuilderImpl) Search(fn func(api.Case)) api.Column {
	cs := &CaseSearchImpl{}
	fn(cs)

	col := dml.Column{From: cs.Build()}
	return &ColumnImpl{Column: col}

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
