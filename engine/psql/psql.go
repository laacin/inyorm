package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/laacin/inyorm"
)

func Open(ctx context.Context, dsn string) *inyorm.Engine {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return &inyorm.Engine{Err: err}
	}

	return &inyorm.Engine{
		Dialect: &PsqlDialect{},
		Driver:  &PsqlDriver{Conn: conn},
	}
}
