package params

import (
	"fmt"
	"maps"
	"strconv"

	"github.com/laacin/inyorm/internal/core"
)

type Store struct {
	ids      []string
	pendings []string
	objects  []objInfo

	store map[string]any
	seens map[string]struct{}

	errs []error
}

func New() core.ParamStore {
	return &Store{
		store: map[string]any{},
		seens: map[string]struct{}{},
	}
}

// Direct store

func (s *Store) Store(v any) {
	id := s.rand()

	s.ids = append(s.ids, id)
	s.store[id] = v
	s.seens[id] = struct{}{}
}

// Lazy store

func (s *Store) Lazy(id string) {
	if id == "" {
		id = s.rand()

		s.ids = append(s.ids, id)
		s.pendings = append(s.pendings, id)
		s.seens[id] = struct{}{}
		return
	}

	if _, exists := s.seens[id]; exists {
		s.pushErr("param conflict: %s already exists", id)
		return
	}

	s.seens[id] = struct{}{}
	s.ids = append(s.ids, id)
}

func (s *Store) Fill(id string, v any) {
	if id == "" {
		if len(s.pendings) < 1 {
			s.pushErr("param overflow: unexpected lazy value")
			return
		}

		id = s.pendings[0]
		s.pendings = s.pendings[1:]

		s.store[id] = v
		return
	}

	if _, exists := s.seens[id]; !exists {
		s.pushErr("param overflow: parameter %s is not declared", id)
		return
	}

	if _, exists := s.store[id]; exists {
		s.pushErr("param duplicate: %s is already assigned", id)
		return
	}

	s.store[id] = v
}

// Lazy object store

func (s *Store) LazyObject(cols []string) {
	baseId := s.rand()

	objIds := make([]string, len(cols))
	for i, col := range cols {
		id := baseId + col

		if _, exists := s.seens[id]; exists {
			s.pushErr("param conflict: %s already exists", id)
			return
		}

		objIds[i] = id
		s.ids = append(s.ids, id)
		s.seens[id] = struct{}{}
	}

	s.objects = append(s.objects, objInfo{
		ids:  objIds,
		cols: cols,
	})
}

func (s *Store) FillObject(fn func(cols []string) []any) {
	if len(s.objects) < 1 {
		s.pushErr("param overflow: unexpected lazy object")
		return
	}

	objInfo := s.objects[0]
	s.objects = s.objects[1:]

	vals := fn(objInfo.cols)
	if len(vals) != len(objInfo.ids) {
		s.pushErr(
			"param object mismatch: expected %d values, got %d",
			len(objInfo.ids),
			len(vals),
		)
		return
	}

	for i, id := range objInfo.ids {
		s.store[id] = vals[i]
	}
}

// Extras

// idx = 0 is the last inserted, idx = 1 is the previous one
func (s *Store) LastIndex(idx int) core.ParamIndex {
	num := len(s.ids)
	if idx >= len(s.ids) {
		return core.ParamIndex{}
	}

	return core.ParamIndex{
		ID:  s.ids[num-idx-1],
		Num: num - idx,
	}
}

func (s *Store) Values() ([]any, error) {
	for _, err := range s.errs {
		if err != nil {
			return nil, err
		}
	}

	vals := make([]any, len(s.ids))
	for i, id := range s.ids {
		val, ok := s.store[id]
		if !ok {
			return nil, fmt.Errorf("param error: %s parameter not found", id)
		}

		vals[i] = val
	}

	return vals, nil
}

func (s *Store) Clone() core.ParamStore {
	clone := &Store{}

	clone.ids = append([]string(nil), s.ids...)
	clone.pendings = append([]string(nil), s.pendings...)
	clone.errs = append([]error(nil), s.errs...)

	clone.objects = make([]objInfo, len(s.objects))
	for i, obj := range s.objects {
		clone.objects[i] = objInfo{
			ids:  append([]string(nil), obj.ids...),
			cols: append([]string(nil), obj.cols...),
		}
	}

	clone.store = make(map[string]any, len(s.store))
	clone.seens = make(map[string]struct{}, len(s.seens))
	maps.Copy(clone.store, s.store)
	maps.Copy(clone.seens, s.seens)

	return clone
}

// helpers

func (s *Store) rand() string {
	id := "_" + strconv.FormatUint(uint64(len(s.ids)+1), 36)
	return id
}

func (s *Store) pushErr(msg string, vals ...any) {
	s.errs = append(s.errs, fmt.Errorf(msg, vals...))
}

type objInfo struct {
	ids  []string
	cols []string
}
