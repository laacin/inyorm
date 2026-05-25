package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

type WriterImpl struct {
	sb    strings.Builder
	dial  expr.ExprWriter
	alias *AliasStore
}

func New(dial expr.ExprWriter, useAliases bool) *WriterImpl {
	if useAliases {
		return &WriterImpl{dial: dial, alias: &AliasStore{}}
	}
	return &WriterImpl{dial: dial}
}

// --- Writer

func (w *WriterImpl) Write(v string) {
	w.sb.WriteString(v)
}

func (w *WriterImpl) Char(v byte) {
	w.sb.WriteByte(v)
}

func (w *WriterImpl) Wrap(fn func(string, core.Writer)) {
	current := w.sb.String()
	w.sb.Reset()
	fn(current, w)
}

func (w *WriterImpl) Value(v any, mode core.WritingMode) {
	expr.NormalizeExpr(v).Render(w, w.dial, mode)
}

func (w *WriterImpl) GetRef(ref string) (byte, bool) {
	if ref == "" {
		return 0, false
	}
	return w.alias.Get(ref)
}

// --- Internal writer

func (w *WriterImpl) SetRef(ref string) {
	if ref == "" {
		return
	}
	w.alias.Set(ref)
}

func (w *WriterImpl) New() core.InternalWriter {
	return &WriterImpl{
		dial:  w.dial,
		alias: w.alias,
	}
}

func (w *WriterImpl) ToString() string {
	return w.sb.String()
}

func (w *WriterImpl) Reset() {
	w.sb.Reset()
}
