package writer

import "errors"

type Placeholder struct {
	count  int
	values []any
}

func (p *Placeholder) Store(v any) {
	p.values = append(p.values, v)
	p.count++
}

func (p *Placeholder) JustCount() {
	p.count++
}

func (p *Placeholder) Validate() error {
	if p.count != len(p.values) {
		return errors.New("mismatched number of parameters and values")
	}
	return nil
}

func (p *Placeholder) GetCount() int {
	return p.count
}

func (p *Placeholder) GetValues() []any {
	return p.values
}

