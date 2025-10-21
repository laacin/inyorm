package internal

import (
	"strconv"
	"strings"
)

const (
	Simple = iota
	Numbered
)

type PlaceholderGen struct {
	Kind   int
	count  int
	values []any
}

func (ph *PlaceholderGen) Write(sb *strings.Builder, values ...any) {
	ph.values = append(ph.values, values...)
	num := len(values)

	if num > 1 {
		sb.WriteByte('(')
	}
	for i := range num {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ph.write())
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
