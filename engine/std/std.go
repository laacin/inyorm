package std

import (
	"database/sql"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/std/std_dialect"
	"github.com/laacin/inyorm/engine/std/std_driver"
)

// Expose dialect for the others dialects
type Dialect = std_dialect.Dialect

func FromInstance(db *sql.DB) *inyorm.Engine {
	return &inyorm.Engine{
		Driver:  &std_driver.Driver{Instance: db},
		Dialect: &std_dialect.Dialect{},
	}
}

func JustDialect() *inyorm.Engine {
	return &inyorm.Engine{
		Dialect: &std_dialect.Dialect{},
	}
}
