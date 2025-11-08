package writer

import (
	"strconv"

	"github.com/laacin/inyorm/internal/core"
)

type Placeholder struct {
	dialect string
	count   int
	values  []any
}

func (ph *Placeholder) next(value any) string {
	ph.count++
	ph.values = append(ph.values, value)

	switch ph.dialect {
	case core.Postgres:
		return "$" + strconv.Itoa(ph.count)
	default:
		return "?"
	}
}
