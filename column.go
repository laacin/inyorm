package inyorm

import (
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

type ColumnImpl struct {
	sb    strings.Builder
	table string
}

func (c *ColumnImpl) Set(col string, from ...string) string {
	c.sb.Reset()

	if len(from) > 0 {
		c.sb.WriteString(from[0])
	} else {
		c.sb.WriteString(c.table)
	}

	c.sb.WriteByte('.')
	c.sb.WriteString(col)

	return stmt.SetColumn(c.sb.String())
}

func (c *ColumnImpl) NewInt(as string, fn func(cb *ColumnBuilder[int]))

func Tst() {
	var c ColumnImpl

	c.NewInt("age", func(cb *ColumnBuilder[int]) {
		cb.Concat("asd", " ", "da")
		c.NewInt("real", func(cb *ColumnBuilder[int]) {
			cb.OnResult(func(result int) {
				result + 10
			})
		})
	})
}
