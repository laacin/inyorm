package mapper

import "errors"

var (
	ErrInvalidSchema = errors.New("invalid schema type")
	ErrInvalidColumn = errors.New("invalid column name")

	ErrPtrExpected   = errors.New("expected pointer")
	ErrValueExpected = errors.New("expected value")

	ErrColumnMismatch = errors.New("column mismatch in slice")
	ErrEmptySlice     = errors.New("empty slice")

	ErrUnexpectedType = errors.New("unexpected type")
)
