package mapper

import "errors"

var (
	ErrPtrExpected   = errors.New("expected pointer")
	ErrValueExpected = errors.New("expected value")

	ErrColumnMismatch = errors.New("column mismatch in slice")
	ErrEmptySlice     = errors.New("empty slice")

	ErrUnexpectedType = errors.New("unexpected type")
)
