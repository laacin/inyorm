package mapper

import "errors"

var (
	ErrSinglePrimitiveColumn  = errors.New("expected exactly one column for primitive type")
	ErrNilPointer             = errors.New("value must be a non-nil pointer")
	ErrInvalidType            = errors.New("invalid type")
	ErrMixedSliceElementTypes = errors.New("slice elements must all be of the same struct type")
	ErrEmptyValues            = errors.New("cannot call WithValues with an empty slice or nil value")

	ErrExpectedStruct        = errors.New("expected struct")
	ErrExpectedSlice         = errors.New("expected slice of structs")
	ErrExpectedPrimitiveType = errors.New("expected primitive type")
	ErrExpectedPointer       = errors.New("expected pointer")
)
