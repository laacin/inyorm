package std_driver

import (
	"context"

	"github.com/laacin/inyorm/internal/ir/driver"
)

func (d *Driver) Exec(ctx context.Context, query string, args ...any) error {
	_, err := d.Instance.ExecContext(ctx, query, args...)
	return err
}

func (d *Driver) Query(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	return d.Instance.QueryContext(ctx, query, args...)
}
