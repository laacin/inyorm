package std_driver

import "database/sql"

type Driver struct {
	Instance *sql.DB
}

func (d *Driver) Close() error {
	return d.Instance.Close()
}
