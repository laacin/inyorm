package clause

import "strconv"

const (
	Simple = iota
	Numbered
)

type PlaceholderGen struct {
	Kind  int
	Count int
}

func (ph *PlaceholderGen) Next() string {
	ph.Count++

	switch ph.Kind {
	case Simple:
		return "?"

	case Numbered:
		return "$" + strconv.Itoa(ph.Count)

	default:
		return "?"
	}
}

func (ph *PlaceholderGen) GetCount() int {
	return ph.Count
}
