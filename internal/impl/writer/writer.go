package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/ir/expr"
)

type WriterImpl struct {
	sb      strings.Builder
	Syntax  expr.ExprWriter
	Aliases *AliasStore
	Params  *ParamStore
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
	Normalize(v).Build(w, w.Syntax, mode)
}

func (w *WriterImpl) ValueCount() int {
	return w.Params.count
}

func (w *WriterImpl) GetRef(ref string) (byte, bool) {
	if ref == "" {
		return 0, false
	}
	return w.Aliases.Get(ref)
}

func (w *WriterImpl) New() core.InternalWriter {
	return &WriterImpl{
		Syntax:  w.Syntax,
		Aliases: w.Aliases,
		Params:  w.Params,
	}
}

func (w *WriterImpl) ToString() string {
	return w.sb.String()
}

func (w *WriterImpl) Reset() {
	w.sb.Reset()
}

// --- Internal writer

func (w *WriterImpl) PushValue(v any) {
	w.Params.Store(v)
}

func (w *WriterImpl) IncValueCount() {
	w.Params.JustCount()
}

func (w *WriterImpl) SetRef(ref string) {
	if ref == "" {
		return
	}
	w.Aliases.Set(ref)
}
