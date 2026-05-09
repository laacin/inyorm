package std

import (
	"database/sql"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/std/dialect"
	"github.com/laacin/inyorm/engine/std/driver"
)

func FromInstance(db *sql.DB) *inyorm.Engine {
	return &inyorm.Engine{
		Driver:  &driver.StdDriver{Instance: db},
		Dialect: &dialect.DialectStd{},
	}
}

func JustDialect() *inyorm.Engine {
	return &inyorm.Engine{
		Dialect: &dialect.DialectStd{},
	}
}
