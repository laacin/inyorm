package query

import "github.com/laacin/inyorm/internal/core"

type Query interface {
	Build(*Tools) error
	Render(core.InternalWriter) error
}
