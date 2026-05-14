package dml

type Statement struct {
	Query  string
	Values []any
}
