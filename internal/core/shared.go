package core

import "context"

// Useful types
type (
	WriterFunc     = func(Writer)
	LazyVal[T any] = func() T
)

// Helpers
func OptionalCtx(candidate []context.Context) context.Context {
	if len(candidate) > 0 && candidate[0] != nil {
		return candidate[0]
	}
	return context.Background()
}
