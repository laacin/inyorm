package execution

import (
	"context"
	"fmt"

	"github.com/laacin/inyorm/internal/entity"
	"github.com/laacin/inyorm/internal/mapper"
)

func run(
	ctx context.Context,
	db entity.Driver,
	query string,
	args []any,
) error {
	if err := db.Exec(ctx, query, args...); err != nil {
		return errSQL(err)
	}
	return nil
}

func scan(
	ctx context.Context,
	db entity.Driver,
	tag string,
	query string,
	args []any,
	scan any,
) error {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return errSQL(err)
	}

	return mapper.Scan(rows, tag, scan)
}

// TODO:

// func runPrep(
// 	ctx context.Context,
// 	stmt *sql.Stmt,
// 	args []any,
// ) error {
// 	_, err := stmt.ExecContext(ctx, args...)
// 	if err != nil {
// 		return errSQL(err)
// 	}
// 	return nil
// }
//
// func scanPrep(
// 	ctx context.Context,
// 	stmt *sql.Stmt,
// 	tag string,
// 	vals []any,
// 	scan any,
// ) error {
// 	rows, err := stmt.QueryContext(ctx, vals...)
// 	if err != nil {
// 		return errSQL(err)
// 	}
//
// 	return mapper.Scan(rows, tag, scan)
// }

func errSQL(err error) error {
	return fmt.Errorf("SQL Error: %w", err)
}
