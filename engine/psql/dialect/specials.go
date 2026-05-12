package dialect

import (
	"strconv"

	"github.com/laacin/inyorm/internal/entity/core"
)

func (dial *PsqlDialect) WritePlaceholder(w core.Writer, count int) {
	w.Char('$')
	w.Write(strconv.Itoa(count))
}
