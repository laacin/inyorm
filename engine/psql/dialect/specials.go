package dialect

import (
	"strconv"

	"github.com/laacin/inyorm/internal/entity"
)

func (dial *PsqlDialect) WritePlaceholder(w entity.Writer, count int) {
	w.Char('$')
	w.Write(strconv.Itoa(count))
}
