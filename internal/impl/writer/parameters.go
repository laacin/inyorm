package writer

import "errors"

type ParamStore struct {
	count  int
	values []any
}

func (p *ParamStore) Store(v any) {
	p.values = append(p.values, v)
}

func (p *ParamStore) JustCount() {
	p.count++
}

func (p *ParamStore) Values() []any {
	return p.values
}

func (p *ParamStore) Validate() error {
	if p.count != len(p.values) {
		return errors.New("mismatched number of parameters and values")
	}
	return nil
}

// type ParamStore struct {
// 	ids   []string
// 	store map[string]any
// }
//
// func (p *ParamStore) Store(v any) string {
// 	id := strconv.Itoa(len(p.ids) + 1)
// 	p.ids = append(p.ids, id)
// 	p.store[id] = v
// 	return id
// }
//
// func (p *ParamStore) FillByID(id string, v any) {
// 	p.store[id] = v
// }
//
// func (p *ParamStore) FillByObj(obj any) {
// 	cols := mapper.ReadColumns([]any{obj})
// 	vals, _ := mapper.ReadValues(cols, obj)
//
// 	if vals.Rows != 1 {
// 		return
// 	}
//
// 	for i := range cols {
// 		p.store[cols[i]] = vals.Args[i]
// 	}
// }
//
// func (p *ParamStore) LazyByID(id string) {
// 	if _, exists := p.store[id]; !exists {
// 		p.ids = append(p.ids, id)
// 	}
// }
//
// func (p *ParamStore) LazyByObj(obj any) {
// 	cols := mapper.ReadColumns([]any{obj})
//
// 	for _, col := range cols {
// 		if _, exists := p.store[col]; !exists {
// 			p.ids = append(p.ids, col)
// 		}
// 	}
// }
//
// // --- Getter
// func (p *ParamStore) Values() ([]any, error) {
// 	vals := make([]any, len(p.ids))
//
// 	for i, id := range p.ids {
// 		v, ok := p.store[id]
// 		if !ok {
// 			return nil, fmt.Errorf("missing parameter: %s", id)
// 		}
//
// 		vals[i] = v
// 	}
//
// 	return vals, nil
// }
//
// func (p *ParamStore) ValueCount() int {
// 	return len(p.ids)
// }
