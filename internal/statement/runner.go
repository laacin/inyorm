package statement

import (
	"context"
	"errors"

	"github.com/laacin/inyorm/internal/core/mapper"
)

var (
	errExecNoDriver = errors.New("missing driver")
)

func Raw(stmt *Statement) (string, []any, error) {
	result, err := stmt.qc.Compile()
	if err != nil {
		return "", nil, err
	}

	vals, err := result.Params.Values()
	if err != nil {
		return "", nil, err
	}

	return result.QueryString, vals, nil
}

func Run(stmt *Statement, ctx context.Context) error {
	if stmt.driver == nil {
		return errExecNoDriver
	}

	query, vals, err := Raw(stmt)
	if stmt.bind == nil {
		return stmt.driver.Exec(ctx, query, vals...)
	}

	rows, err := stmt.driver.Query(ctx, query, vals...)
	if err != nil {
		return err
	}

	return mapper.New().Bind(rows, stmt.bind)
}
