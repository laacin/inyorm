package core

// Useful types
type (
	WriterFunc     = func(Writer)
	LazyVal[T any] = func() T
)
