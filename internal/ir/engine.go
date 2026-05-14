package ir

import (
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
	dml.Syntax
	expr.Syntax
}
