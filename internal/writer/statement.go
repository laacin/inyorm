package writer

import "github.com/laacin/inyorm/internal/core"

type Statement struct {
	Ph    Placeholder
	Alias Alias
}

func (stmt *Statement) Writer() core.Writer {
	return &Writer{ph: &stmt.Ph, aliases: &stmt.Alias}
}
