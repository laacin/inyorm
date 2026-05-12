package driver

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/laacin/inyorm/internal/entity/driver"
)

func (d *PsqlDriver) Exec(ctx context.Context, query string, args ...any) error {
	_, err := d.Conn.Exec(ctx, query, args...)
	return err
}

func (d *PsqlDriver) Query(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	rows, err := d.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &PsqlRows{rows}, nil
}

// --- Dependencies

type PsqlRows struct {
	pgx.Rows
}

func (r *PsqlRows) Columns() ([]string, error) {
	result := r.Rows.FieldDescriptions()

	cols := make([]string, len(result))
	for i, rslt := range result {
		cols[i] = rslt.Name
	}

	return cols, nil
}
