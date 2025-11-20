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

func (ph *Placeholder) withValue(value any) string {
	ph.values = append(ph.values, value)
	return ph.write()
}

func (ph *Placeholder) write() string {
	ph.count++

	switch ph.dialect {
	case core.Postgres:
		return "$" + strconv.Itoa(ph.count)
	default:
		return "?"
	}
}
