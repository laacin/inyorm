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

func (a *AliasStore) Get(table string) byte {
	if a.list == nil {
		a.list = make(map[string]byte)
	}

	if al, exists := a.list[table]; exists {
		return al
	}

	alias := abc[a.count]
	a.list[table] = alias
	a.count++
	return alias
}
