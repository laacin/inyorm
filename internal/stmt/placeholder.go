package stmt

import (
	"strconv"
	"strings"
)

const (
	Simple = iota
	Numbered
)

type PlaceholderGen struct {
	Stringify bool
	Kind      int
	count     int
	values    []any
}

func (ph *PlaceholderGen) Write(sb *strings.Builder, values ...any) {
	if !ph.Stringify {
		ph.values = append(ph.values, values...)
	}
	num := len(values)

	if num > 1 {
		sb.WriteByte('(')
	}
	for i := range num {
		if i > 0 {
			sb.WriteString(", ")
		}
		if !ph.Stringify {
			sb.WriteString(ph.write())
		} else {
			str := Stringify(values[i])
			sb.WriteString(str)
		}
	}
	if num > 1 {
		sb.WriteByte(')')
	}
}

func (ph *PlaceholderGen) Values() []any { return ph.values }

func (ph *PlaceholderGen) write() string {
	ph.count++

	switch ph.Kind {
	case Simple:
		return "?"

	case Numbered:
		return "$" + strconv.Itoa(ph.count)

	default:
		return "?"
	}
}
