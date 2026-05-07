package writer

var abc = [...]byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

type AliasStore struct {
	count int
	list  map[string]byte
}

func (a *AliasStore) Get(ref string) (byte, bool) {
	if a == nil || a.list == nil {
		return 0, false
	}

	if al, exists := a.list[ref]; exists {
		return al, true
	}

	return 0, false
}

func (a *AliasStore) Set(ref string) {
	if a == nil {
		return
	}

	if a.list == nil {
		a.list = make(map[string]byte)
	}

	if _, exists := a.list[ref]; exists {
		return
	}

	alias := abc[a.count]
	a.list[ref] = alias
	a.count++
}
