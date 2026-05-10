package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/engine/psql/dialect"
	"github.com/laacin/inyorm/engine/psql/driver"
)

type PsqlDialect = dialect.PsqlDialect

func Open(ctx context.Context, dsn string) *inyorm.Engine {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return &inyorm.Engine{Err: err}
	}

	return &inyorm.Engine{
		Dialect: &dialect.PsqlDialect{},
		Driver:  &driver.PsqlDriver{Conn: conn},
	}
}
