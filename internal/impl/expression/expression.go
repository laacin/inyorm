package expression

import (
	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/entity/api"
)

type ExpressionImpl struct {
	MainRef string
	Dialect entity.ValueWriter
}

func (expr *ExpressionImpl) Table(name string) api.Table {
	return &entity.Table{Value: name}
}

func (expr *ExpressionImpl) Col(name string, ref ...string) api.Column {
	col := entity.Column{Ref: getLast(expr.MainRef, ref), Name: name}
	return &ColumnImpl{Column: col}
}

func (expr *ExpressionImpl) All(ref ...string) api.Column {
	col := entity.Column{Ref: getLast(expr.MainRef, ref), From: &entity.Wildcard{}}
	return &ColumnImpl{Column: col}
}

func (expr *ExpressionImpl) Param(value ...any) api.Parameter {
	return &entity.Parameter{Store: len(value) > 0, Value: getLast(nil, value)}
}

func (expr *ExpressionImpl) Cond(ident any) api.Condition {
	cond := &ConditionImpl{}
	return cond.Start(ident)
}

func (expr *ExpressionImpl) Concat(values ...any) api.Column {
	col := entity.Column{From: &entity.Concat{Values: values}}
	return &ColumnImpl{Column: col}
}

func (expr *ExpressionImpl) Switch(cond any, fn func(api.Case)) api.Column {
	cs := &CaseSwitchImpl{}
	fn(cs)

	col := entity.Column{From: cs.Build()}
	return &ColumnImpl{Column: col}
}

func (expr *ExpressionImpl) Search(fn func(api.Case)) api.Column {
	cs := &CaseSearchImpl{}
	fn(cs)

	col := entity.Column{From: cs.Build()}
	return &ColumnImpl{Column: col}

}

// --- Helpers
func getLast[T any](prev T, candidate []T) T {
	if len(candidate) > 0 {
		return candidate[0]
	}
	return prev
}
