package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/entity"
)

type WriterImpl struct {
	sb      strings.Builder
	vw      entity.ValueWriter
	aliases *AliasStore
	params  *ParamStore
}

func (w *WriterImpl) Write(v string) {
	w.sb.WriteString(v)
}

func (w *WriterImpl) Char(v byte) {
	w.sb.WriteByte(v)
}

func (w *WriterImpl) Value(v any, mode entity.WritingMode) {
	Normalize(v).Write(w, w.vw, mode)
}

func (w *WriterImpl) ValueCount() int {
	return w.params.count
}

func (w *WriterImpl) GetRef(table string) (byte, bool) {
	if w.aliases != nil {
		return w.aliases.Get(table), true
	}
	return ' ', false
}

func (w *WriterImpl) New() entity.Writer {
	return &WriterImpl{vw: w.vw}
}

func (w *WriterImpl) Result() string {
	return w.sb.String()
}

func (w *WriterImpl) Reset() {
	w.sb.Reset()
}
