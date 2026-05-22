package sqlite

import (
	"database/sql"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/std"
	_ "github.com/mattn/go-sqlite3"
)

func Open(dest string) *inyorm.Engine {
	db, err := sql.Open("sqlite3", dest)
	if err != nil {
		return &inyorm.Engine{Err: err}
	}
	return std.FromInstance(db)
}
