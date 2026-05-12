package dialect

import (
	"strconv"

	"github.com/laacin/inyorm/internal/entity/core"
)

var quote byte = "'"[0]

// --- Literals
func (dial *StdDialect) WriteString(w core.Writer, v string) {
	w.Char(quote)
	w.Write(v)
	w.Char(quote)
}

func (dial *StdDialect) WriteNumber(w core.Writer, v int) {
	r := strconv.Itoa(v)
	w.Write(r)
}

func (dial *StdDialect) WriteFloat(w core.Writer, v float64) {
	r := strconv.FormatFloat(float64(v), 'f', -1, 32)
	w.Write(r)
}

func (dial *StdDialect) WriteBool(w core.Writer, v bool) {
	if v {
		w.Char('1')
		return
	}
	w.Char('0')
}

func (dial *StdDialect) WriteNull(w core.Writer) {
	w.Write("NULL")
}

func (dial *StdDialect) WriteWildcard(w core.Writer) {
	w.Char('*')
}
