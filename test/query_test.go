package test

import (
	"testing"

	"github.com/laacin/inyorm"
)

// -- Schema
type UserColumn struct {
	ID        inyorm.Column
	Firstname inyorm.Column
	Lastname  inyorm.Column
	Age       inyorm.Column
	Banned    inyorm.Column
}

var UserTable = inyorm.TableFor[UserColumn]{
	Table: inyorm.Table{
		Name:       "users",
		PrimaryKey: "id",
		Foreigns:   nil,
	},
	Column: UserColumn{
		ID:        "id",
		Firstname: "firstname",
		Lastname:  "lastname",
		Age:       "age",
		Banned:    "banned",
	},
}

func TestQuery(t *testing.T) {
	q := inyorm.NewQuery(UserTable.Table)
	q.Where()
}
