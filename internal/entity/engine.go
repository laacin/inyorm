package entity

import (
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/driver"
	"github.com/laacin/inyorm/internal/entity/expr"
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
