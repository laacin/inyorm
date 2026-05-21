package std_driver

import (
	"context"
	"database/sql"

	"github.com/laacin/inyorm/internal/core"
)

type Driver struct {
	Instance *sql.DB
}

type SQLTx struct {
	tx *sql.Tx
}

type SQLPrepared struct {
	stmt *sql.Stmt
}

type SQLRows struct {
	rows *sql.Rows
}

func (d *Driver) Close() error {
	return d.Instance.Close()
}

func (d *Driver) BeginTx(ctx context.Context) core.Transaction {
	tx, err := d.Instance.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	return &SQLTx{
		tx: tx,
	}
}

func (d *Driver) Exec(
	ctx context.Context,
	query string,
	args ...any,
) error {
	_, err := d.Instance.ExecContext(ctx, query, args...)
	return err
}

func (d *Driver) Query(
	ctx context.Context,
	query string,
	args ...any,
) (core.Rows, error) {
	rows, err := d.Instance.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &SQLRows{
		rows: rows,
	}, nil
}

func (d *Driver) Prepare(
	ctx context.Context,
	query string,
) (core.Prepared, error) {
	stmt, err := d.Instance.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &SQLPrepared{
		stmt: stmt,
	}, nil
}

func (t *SQLTx) Commit() error {
	return t.tx.Commit()
}

func (t *SQLTx) Rollback() error {
	return t.tx.Rollback()
}

func (t *SQLTx) Exec(
	ctx context.Context,
	query string,
	args ...any,
) error {
	_, err := t.tx.ExecContext(ctx, query, args...)
	return err
}

func (t *SQLTx) Query(
	ctx context.Context,
	query string,
	args ...any,
) (core.Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &SQLRows{
		rows: rows,
	}, nil
}

func (t *SQLTx) Prepare(
	ctx context.Context,
	query string,
) (core.Prepared, error) {
	stmt, err := t.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &SQLPrepared{
		stmt: stmt,
	}, nil
}

func (p *SQLPrepared) Exec(
	ctx context.Context,
	args ...any,
) error {
	_, err := p.stmt.ExecContext(ctx, args...)
	return err
}

func (p *SQLPrepared) Query(
	ctx context.Context,
	args ...any,
) (core.Rows, error) {
	rows, err := p.stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	return &SQLRows{
		rows: rows,
	}, nil
}

func (p *SQLPrepared) Close() error {
	return p.stmt.Close()
}

func (r *SQLRows) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *SQLRows) Next() bool {
	return r.rows.Next()
}

func (r *SQLRows) Scan(args ...any) error {
	return r.rows.Scan(args...)
}

func (r *SQLRows) Err() error {
	return r.rows.Err()
}

func (r *SQLRows) Close() error {
	return r.rows.Close()
}
