package driver

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PsqlDriver struct {
	Conn *pgx.Conn
}

func (d *PsqlDriver) Close() error {
	return d.Conn.Close(context.Background())
}
