package test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/laacin/inyorm"
)

// --- COMMON HELPERS

func Start(eng *inyorm.Engine) *inyorm.DB {
	db, err := inyorm.New(eng)
	if err != nil {
		panic(err)
	}

	return db
}

func End(db *inyorm.DB, cb ...func() error) {
	if err := db.Close(); err != nil {
		panic(err)
	}

	if len(cb) > 0 {
		if err := cb[0](); err != nil {
			panic(err)
		}
	}
}

// --- SQLITE

var tmpSqliteFilePath = func() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, "data_test.db")
	if !strings.HasSuffix(path, "inyorm/test/data_test.db") {
		panic("Wrong path")
	}

	return path
}()

func deleteSqliteFile() error { return os.Remove(tmpSqliteFilePath) }
