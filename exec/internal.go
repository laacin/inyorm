package exec

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/laacin/inyorm/internal/mapper"
)

type instance interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func run(
	ctx context.Context,
	db instance,
	query string,
	args []any,
) error {
	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return errSQL(err)
	}
	return nil
}

func scan(
	ctx context.Context,
	db instance,
	tag string,
	query string,
	args []any,
	binder any,
) error {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return errSQL(err)
	}
	return mapper.Scan(rows, tag, binder)
}

func runPrep(
	ctx context.Context,
	stmt *sql.Stmt,
	args []any,
) error {
	_, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return errSQL(err)
	}
	return nil
}

func scanPrep(
	ctx context.Context,
	stmt *sql.Stmt,
	tag string,
	vals []any,
	binder any,
) error {
	rows, err := stmt.QueryContext(ctx, vals...)
	if err != nil {
		return errSQL(err)
	}
	return mapper.Scan(rows, tag, binder)
}

func errSQL(err error) error {
	return fmt.Errorf("SQL Error: %w", err)
}
