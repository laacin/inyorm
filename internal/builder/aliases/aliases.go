package aliases

import "github.com/laacin/inyorm/internal/core"

var abc = [...]byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

type AliasStore struct {
	count   int
	list    map[string]byte
	enabled bool
	main    string
}

func New() *AliasStore {
	return &AliasStore{}
}

func (a *AliasStore) Enable() {
	a.enabled = true
}

func (a *AliasStore) Get(ref string) core.Reference {
	return a.getIfIsPossible(ref)
}

func (a *AliasStore) Set(ref string) {
	if ref == a.main {
		return
	}

	if a.store(ref, a.count) {
		a.count++
	}
}

func (a *AliasStore) SetMain(ref string) {
	a.main = ref
	a.store(ref, 0)
}

func (a *AliasStore) GetMain() core.Reference {
	return a.getIfIsPossible(a.main)
}

func (a *AliasStore) getIfIsPossible(ref string) core.Reference {
	if a.list == nil {
		return core.Reference{}
	}

	if al, ok := a.list[ref]; ok {
		return core.Reference{
			Ref:    al,
			Enable: a.enabled,
		}
	}

	return core.Reference{}
}

func (a *AliasStore) store(ref string, idx int) bool {
	if a.list == nil {
		a.list = make(map[string]byte)
	}

	if _, exists := a.list[ref]; exists {
		return false
	}

	if ref == a.main {
		a.list[ref] = abc[0]
		return true
	}

	a.list[ref] = abc[idx]
	return true
}
