package internal

import "strconv"

type Placeholder struct {
	dialect string
	count   int
	values  []any
}

func (ph *Placeholder) next(value any) string {
	ph.count++
	ph.values = append(ph.values, value)

	switch ph.dialect {
	case psql:
		return "$" + strconv.Itoa(ph.count)
	default:
		return "?"
	}
}
