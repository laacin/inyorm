package stmt

import (
	"strconv"
	"strings"
)

const (
	SimplePh = iota
	NumberedPh
)

type PlaceholderGen struct {
	StringMode bool
	Kind       int
	count      int
	values     []any
}

func (ph *PlaceholderGen) Write(sb *strings.Builder, values ...any) {
	if !ph.StringMode {
		ph.values = append(ph.values, values...)
	}
	num := len(values)

	if num > 1 {
		sb.WriteByte('(')
		defer sb.WriteByte(')')
	}

	for i := range num {
		if i > 0 {
			sb.WriteString(", ")
		}

		if !ph.StringMode {
			sb.WriteString(ph.write())
		} else {
			str := Stringify(values[i])
			sb.WriteString(str)
		}
	}
}

func (ph *PlaceholderGen) Values() []any { return ph.values }

func (ph *PlaceholderGen) write() string {
	ph.count++

	switch ph.Kind {
	case SimplePh:
		return "?"

	case NumberedPh:
		return "$" + strconv.Itoa(ph.count)

	default:
		return "?"
	}
}
