package params

import (
	"fmt"
	"strconv"

	"github.com/laacin/inyorm/internal/core"
)

type ParamStore struct {
	ids  []string
	lazy []string
	obj  []objInfo

	store map[string]any
	errs  []error
}

func New() *ParamStore { return &ParamStore{} }

func (p *ParamStore) Store(v any) {
	p.initMap()

	id := p.rand()
	p.ids = append(p.ids, id)
	p.store[id] = v
}

func (p *ParamStore) Lazy(id string) {
	if id == "" {
		id = p.rand()

		p.ids = append(p.ids, id)
		p.lazy = append(p.lazy, id)
		return
	}

	if _, exists := p.store[id]; exists {
		p.pushErr("param conflict: %s already exists", id)
		return
	}

	p.ids = append(p.ids, id)
}

func (p *ParamStore) LazyObj(cols []string) {
	baseId := p.rand()

	objIds := make([]string, len(cols))
	for i, col := range cols {
		id := baseId + col
		objIds[i] = id
		p.ids = append(p.ids, id)
	}

	p.obj = append(p.obj, objInfo{
		ids:  objIds,
		cols: cols,
	})
}

func (p *ParamStore) Fill(id string, v any) {
	p.initMap()

	if id != "" {
		if _, exists := p.store[id]; exists {
			p.pushErr("param duplicate: %s is already assigned", id)
			return
		}

		p.store[id] = v
		return
	}

	if len(p.lazy) < 1 {
		p.pushErr("param overflow: unexpected lazy value")
		return
	}

	id = p.lazy[0]
	p.lazy = p.lazy[1:]

	p.store[id] = v
}

func (p *ParamStore) FillObj(fn func(cols []string) []any) {
	p.initMap()

	if len(p.obj) < 1 {
		p.pushErr("param overflow: unexpected lazy object")
		return
	}

	objInfo := p.obj[0]
	p.obj = p.obj[1:]
	vals := fn(objInfo.cols)

	for i, id := range objInfo.ids {
		p.store[id] = vals[i]
	}
}

// idx = 0 is the last inserted, idx = 1 is the previous one
func (p *ParamStore) LastIndex(idx int) core.ParamIndex {
	num := len(p.ids)
	if idx > len(p.ids) {
		return core.ParamIndex{}
	}

	return core.ParamIndex{
		ID:  p.ids[num-idx-1],
		Num: num - idx,
	}
}

func (p *ParamStore) Values() ([]any, error) {
	if p.store == nil {
		return []any{}, nil
	}

	for _, err := range p.errs {
		if err != nil {
			return nil, err
		}
	}

	vals := make([]any, len(p.ids))
	for i, id := range p.ids {
		val, ok := p.store[id]
		if !ok {
			return nil, fmt.Errorf("param error: %s parameter not found", id)
		}

		vals[i] = val
	}

	return vals, nil
}

// --- helpers
func (p *ParamStore) initMap() {
	if p.store == nil {
		p.store = make(map[string]any)
	}
}

func (p *ParamStore) rand() string {
	id := "_" + strconv.FormatUint(uint64(len(p.ids)+1), 36)
	return id
}

func (p *ParamStore) pushErr(msg string, vals ...any) {
	p.errs = append(p.errs, fmt.Errorf(msg, vals...))
}

type objInfo struct {
	ids  []string
	cols []string
}
