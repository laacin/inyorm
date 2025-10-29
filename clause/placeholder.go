package clause

import (
	"strconv"
	"strings"

	"github.com/laacin/inyorm/internal/stmt"
)

const (
	psql = "postgres"
)

type Placeholder struct {
	Dialect string
	count   int
	values  []any
}

func (ph *Placeholder) Write(sb *strings.Builder, values ...any) {
	initialized := ph != nil
	if initialized {
		ph.values = append(ph.values, values...)
	}

	if len(values) > 1 {
		sb.WriteByte('(')
		defer sb.WriteByte(')')
	}

	for i := range len(values) {
		if i > 0 {
			sb.WriteString(", ")
		}

		if initialized {
			sb.WriteString(ph.write())
		} else {
			sb.WriteString(stmt.Stringify(values[i]))
		}
	}
}

func (ph *Placeholder) Values() []any { return ph.values }
func (ph *Placeholder) write() string {
	ph.count++

	switch ph.Dialect {
	case psql:
		return "$" + strconv.Itoa(ph.count)
	default:
		return "?"
	}
}
