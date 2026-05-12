package entity

import (
	"github.com/laacin/inyorm/internal/entity/dml"
	"github.com/laacin/inyorm/internal/entity/driver"
)

type Engine struct {
	DML    dml.Dialect
	Driver driver.Driver
	Err    error
}
