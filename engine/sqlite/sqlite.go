package sqlite

import (
	"database/sql"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/std"
	_ "github.com/mattn/go-sqlite3"
)

func Open(dest string) (*inyorm.Engine, error) {
	db, err := sql.Open("sqlite3", dest)
	if err != nil {
		return nil, err
	}
	return std.FromInstance(db)
}
