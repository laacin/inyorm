package inyorm

type Field string

type CustomField struct {
	writen   *string
	distinct bool
}

type FieldBuilder struct {
	as       string
	distinct string
}

func NewField(as string, fn func(fb *FieldBuilder)) *CustomField

func (fb *FieldBuilder) Op(v Value) *FieldOperation
func (fb *FieldBuilder) Simple(v Value) *CustomField
func (fb *FieldBuilder) Concat(v ...Value) *CustomField

func tst() {
	asd := NewField("asd", func(fb *FieldBuilder) {
		fb.Op("ani").Add(10).Mul(2)
	})
}
