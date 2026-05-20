package api

import "context"

// SelectQuery represents a full SELECT statement
type SelectQuery interface {
	Select
	From
	Join
	Where
	GroupBy
	Having
	OrderBy
	Limit
	Offset
}

// InsertQuery represents a full INSERT statement
type InsertQuery interface {
	Insert
	Into
	Values
}

// UpdateQuery represents a full UPDATE statement
type UpdateQuery interface {
	Update
	Into
	Values
	Where
}

// DeleteQuery represents a full DELETE statement
type DeleteQuery interface {
	Delete
	From
	Where
}

type Statement interface {
	Runner
	Prepare() PrepStatement
	Bind(...any) Statement
}

type PrepStatement interface {
	BindPrep(...any) PrepStatement
	Values(...any) Runner
}

type Runner interface {
	Raw() (string, []any, error)
	Run(...context.Context) error
}
