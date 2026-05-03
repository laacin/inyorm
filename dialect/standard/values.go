package standard

import (
	"strconv"

	"github.com/laacin/inyorm/intr/dialect"
)

func (dial *DialectStandard) String(w dialect.Writer, v string) {
	w.Write("'")
	w.Write(v)
	w.Write("'")
}

func (dial *DialectStandard) Number(w dialect.Writer, v int) {
	w.Write(strconv.Itoa(v))
}

func (dial *DialectStandard) Float(w dialect.Writer, v float64) {
	r := strconv.FormatFloat(float64(v), 'f', -1, 32)
	w.Write(r)
}

func (dial *DialectStandard) Bool(w dialect.Writer, v bool) {
	if v {
		w.Char('1')
	} else {
		w.Char('0')
	}
}

func (dial *DialectStandard) Null(w dialect.Writer) {
	w.Write("NULL")
}

func (dial *DialectStandard) Placeholder(w dialect.Writer, num int) {
	w.Char('?')
}
