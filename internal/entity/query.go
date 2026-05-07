package entity

type Query struct {
	Statement string
	Values    []any
	Errs      []error
}

func (q *Query) Result() (string, []any) {
	return q.Statement, q.Values
}

func (q *Query) FirstErr() error {
	for _, err := range q.Errs {
		if err != nil {
			return err
		}
	}
	return nil
}
