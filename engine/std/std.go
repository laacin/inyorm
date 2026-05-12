package std

import (
	"database/sql"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/std/dialect"
	"github.com/laacin/inyorm/engine/std/driver"
)

// Expose dialect for the others dialects
type StdDialect = dialect.StdDialect

func FromInstance(db *sql.DB) *inyorm.Engine {
	return &inyorm.Engine{
		Driver: &driver.StdDriver{Instance: db},
		DML:    &dialect.StdDialect{},
	}
}

func JustDialect() *inyorm.Engine {
	return &inyorm.Engine{
		DML: &dialect.StdDialect{},
	}
}
