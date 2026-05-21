package sqlite

import (
	"database/sql"

	"github.com/laacin/inyorm/engine/std"
	"github.com/laacin/inyorm/internal/ir"
	_ "github.com/mattn/go-sqlite3"
)

func Open(dest string) *ir.Engine {
	db, err := sql.Open("sqlite3", dest)
	if err != nil {
		return &ir.Engine{Err: err}
	}
	return std.FromInstance(db)
}
