package writer

import "errors"

type ParamStore struct {
	count  int
	values []any
}

func (p *ParamStore) Store(v any) {
	p.values = append(p.values, v)
	p.count++
}

func (p *ParamStore) JustCount() {
	p.count++
}

func (p *ParamStore) Validate() error {
	if p.count != len(p.values) {
		return errors.New("mismatched number of parameters and values")
	}
	return nil
}
