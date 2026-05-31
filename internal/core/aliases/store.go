package aliases

import "github.com/laacin/inyorm/internal/core"

var abc = [...]byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

type Store struct {
	count   int
	list    map[string]byte
	enabled bool
	main    string
}

func New() *Store {
	return &Store{}
}

func (s *Store) Enable() {
	s.enabled = true
}

func (s *Store) Get(ref string) core.Reference {
	return s.getIfIsPossible(ref)
}

func (s *Store) Set(ref string) {
	if ref == s.main {
		return
	}

	if s.store(ref, s.count) {
		s.count++
	}
}

func (s *Store) SetMain(ref string) {
	s.main = ref
	s.store(ref, 0)
}

func (s *Store) GetMain() core.Reference {
	return s.getIfIsPossible(s.main)
}

func (s *Store) getIfIsPossible(ref string) core.Reference {
	if s.list == nil {
		return core.Reference{}
	}

	if al, ok := s.list[ref]; ok {
		return core.Reference{
			Ref:     al,
			Enabled: s.enabled,
		}
	}

	return core.Reference{}
}

func (s *Store) store(ref string, idx int) bool {
	if s.list == nil {
		s.list = make(map[string]byte)
	}

	if _, exists := s.list[ref]; exists {
		return false
	}

	if ref == s.main {
		s.list[ref] = abc[0]
		return true
	}

	s.list[ref] = abc[idx]
	return true
}
