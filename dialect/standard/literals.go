package standard

import (
	"strconv"

	"github.com/laacin/inyorm/internal/entity"
)

// --- Literals
func (dial *DialectStandard) WriteString(w entity.Writer, v string) {
	quote := "'"[0]

	w.Char(quote)
	w.Write(v)
	w.Char(quote)
}

func (dial *DialectStandard) WriteNumber(w entity.Writer, v int) {
	r := strconv.Itoa(v)
	w.Write(r)
}

func (dial *DialectStandard) WriteFloat(w entity.Writer, v float64) {
	r := strconv.FormatFloat(float64(v), 'f', -1, 32)
	w.Write(r)
}

func (dial *DialectStandard) WriteBool(w entity.Writer, v bool) {
	if v {
		w.Char('1')
	}
	w.Char('0')
}

func (dial *DialectStandard) WriteNull(w entity.Writer) {
	w.Write("NULL")
}

func (dial *DialectStandard) WriteWildcard() string {
	return "*"
}
