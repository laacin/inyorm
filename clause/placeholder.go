package clause

import "strconv"

const (
	Simple = iota
	Numbered
)

type PlaceholderGen struct {
	Kind  int
	count int
}

func (ph *PlaceholderGen) Next() string {
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

func (ph *PlaceholderGen) GetCount() int {
	return ph.count
}
