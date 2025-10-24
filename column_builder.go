package inyorm

import (
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

// --- Field
// type Field string
//
// func (f Field) Set() string { return stmt.SetColumn(string(f)) }

// Aggregations
type Aggregation string

const (
	Count = "COUNT"
	Sum   = "SUM"
	Max   = "MAX"
	Min   = "MIN"
	Avg   = "AVG"
)

type Fieldd string

func (f Fieldd) Set() string

func (f Fieldd) Greater(v any) Fieldd
func (f Fieldd) Less(v any) Fieldd

func (f Fieldd) Add(v any) Fieldd
func (f Fieldd) Sub(v any) Fieldd
func (f Fieldd) Mul(v any) Fieldd
func (f Fieldd) Div(v any) Fieldd
func (f Fieldd) Mod(v any) Fieldd
func (f Fieldd) Get() Fieldd

func (f Fieldd) Substring(start, end int) Fieldd
func (f Fieldd) Upper() Fieldd
func (f Fieldd) Lower() Fieldd
func (f Fieldd) Trim() Fieldd

type ColB struct {
	target Field
}

func (c *ColB) Build(name string, fn func(fb FieldBuilder) Field) Field
func (c *ColB) New(as string, vals ...any) Field

func Tstt() {
	const (
		fname Fieldd = "firstname"
		lname Fieldd = "lastname"
	)

	var c ColB
	var firstname, lastname, age Field
	var price, stock Field

	fullName := c.Build("full_name", func(fb FieldBuilder) Field {
		return fb.Target(firstname.Set(), " ", lastname.Set()).Upper()
	})

	isAdult := c.Build("adult", func(fb FieldBuilder) Field {
		return fb.Search(func(c Case) {
			c.When(age.Greater(17)).Then(true).Else(false)
		})
	})

	total := c.Build("total", func(fb FieldBuilder) Field {
		return price.Mul(stock)
	})
}
