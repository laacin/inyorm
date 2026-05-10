package driver

import "database/sql"

type StdDriver struct {
	Instance *sql.DB
}

func (d *StdDriver) Close() error {
	return d.Instance.Close()
}
