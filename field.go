package inyorm

import (
	"github.com/laacin/inyorm/internal/stmt"
	"strings"
)

type Field string

func (Field) Set() string

type CustomField struct {
	as          string
	targets     []any
	distinct    bool
	selectValue *string
	aggr        Aggregation
	functions   []string
	// logical
	greaterThan any
	lessThan    any
}
type FB struct {
}

func New(as string, fn func(fb *FB) *CustomField) Field

func (fb *FB) Simple(value any) *CustomField
func (fb *FB) Concat(values ...any) *CustomField

// TODO: subquery

func NewWithConcat(as string, values ...any) *CustomField {
	return &CustomField{
		as:      as,
		targets: values,
	}
}

func sttst() {
	var fname, lname, age, prod, qt, price Field
	fdl := New("full_name", func(fb *FB) *CustomField {
		total := (fb.Op(qt).Mul(price)).Mul(100)
		result := fb.Concat(fname, " ", lname, ", Have the Product: ", prod, ", Worth at ", total.End())
		value := fb.Search(func(cs *Case[*ExpressionEnd]) {
			expr := fb.Expr(age).Less(18).Or().Greater(70)
			cs.When(expr).Then(fb.Concat(fname, " ", lname, " is minor age"))
			cs.Else(result)
		})
		return value
	})
}

// Operation
type OpTarget struct{ sb *strings.Builder }

func (cf *FB) Op(v any) *OpTarget
func (op *OpTarget) Add(v any) *OpTarget
func (op *OpTarget) Sub(v any) *OpTarget
func (op *OpTarget) Mul(v any) *OpTarget
func (op *OpTarget) Div(v any) *OpTarget
func (op *OpTarget) Mod(v any) *OpTarget
func (op *OpTarget) End() *CustomField

// conditionals
type Case[T any] struct{}

func (c *Case[T]) When(v T) *CaseResult[T]
func (c *Case[T]) Else(v any)

type CaseResult[T any] struct{}

func (c *CaseResult[T]) Then(any) Case[T]

func (fb *FB) Expr(v any) *Expression

func (fb *FB) Switch(literal any, fn func(cs *Case[any])) *CustomField
func (fb *FB) Search(fn func(cs *Case[*ExpressionEnd])) *CustomField

// Builder
func (cf *CustomField) Use(onSelect ...bool) Field {
	if !(len(onSelect) > 0 && onSelect[0]) {
		return Field(stmt.SetColumn(cf.as))
	}

	if sel := cf.selectValue; sel != nil {
		return Field(*sel)
	}

	var sb strings.Builder
	func() {
		if cf.distinct {
			sb.WriteString("DISTINCT ")
		}

		if cf.aggr != "" {
			sb.WriteString(string(cf.aggr))
			sb.WriteByte('(')
			defer sb.WriteByte(')')
		}

		for _, fn := range cf.functions {
			sb.WriteString(fn)
			sb.WriteByte('(')
			defer sb.WriteByte(')')
		}

		for _, t := range cf.targets {
			sb.WriteString(stmt.Stringify(t))
		}
	}()

	concated := sb.String()
	cf.selectValue = &concated
	return Field(concated)
}
