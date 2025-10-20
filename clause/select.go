package clause

import (
	"errors"
	"strings"
)

type SelectBuilder struct {
	Default string
	Values  []*SelectSimple
}

func (s *SelectBuilder) New(sel string, as ...string) {
	value := &SelectSimple{value: sel}
	if len(as) > 0 && as[0] != "" {
		value.as = as[0]
	}

	s.Values = append(s.Values, value)
}

func (s *SelectBuilder) Build(sb *strings.Builder) error {
	if len(s.Values) == 0 {
		if s.Default == "" {
			return errors.New("no select fields specified and no default select provided")
		}

		s.Values = append(s.Values, &SelectSimple{value: s.Default})
	}

	sb.WriteString("SELECT ")
	for i, sel := range s.Values {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(sel.value)
		if sel.as != "" {
			sb.WriteString(" AS ")
			sb.WriteString(sel.as)
		}
	}

	return nil
}

type SelectSimple struct {
	value string
	as    string
}
