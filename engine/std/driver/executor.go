package driver

import (
	"context"

	"github.com/laacin/inyorm/internal/entity"
)

func (d *StdDriver) Exec(ctx context.Context, query string, args ...any) error {
	_, err := d.Instance.ExecContext(ctx, query, args...)
	return err
}

func (d *StdDriver) Query(ctx context.Context, query string, args ...any) (entity.Rows, error) {
	return d.Instance.QueryContext(ctx, query, args...)
}
