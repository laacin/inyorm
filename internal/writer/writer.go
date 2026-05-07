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

func New(dial entity.Dialect, useAliases bool) *WriterImpl {
	w := &WriterImpl{
		vw:     dial,
		params: &ParamStore{},
	}

	if useAliases {
		w.aliases = &AliasStore{}
	}

	return w
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

func (w *WriterImpl) StoreValue(v any) {
	w.params.Store(v)
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

func (w *WriterImpl) GetValues() ([]any, error) {
	if err := w.params.Validate(); err != nil {
		return nil, err
	}
	return w.params.values, nil
}

func (w *WriterImpl) DefaultAlias(tbl string) {
	w.aliases.Get(tbl)
}
