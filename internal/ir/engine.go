package ir

import (
	"github.com/laacin/inyorm/internal/ir/ddl"
	"github.com/laacin/inyorm/internal/ir/dml"
	"github.com/laacin/inyorm/internal/ir/driver"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type Engine struct {
	Dialect Dialect
	Driver  driver.Driver
	Err     error
}

type Dialect interface {
	ddl.Syntax
	dml.Syntax
	expr.ExprWriter
}
