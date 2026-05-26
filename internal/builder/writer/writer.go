package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/expr"
)

type Writer struct {
	sb   strings.Builder
	dial expr.ExprWriter
}

func New(dial expr.ExprWriter) *Writer {
	return &Writer{dial: dial}
}

// --- Writer

func (w *Writer) Write(v string) {
	w.sb.WriteString(v)
}

func (w *Writer) Char(v byte) {
	w.sb.WriteByte(v)
}

func (w *Writer) Wrap(fn func(string, core.Writer)) {
	current := w.sb.String()
	w.sb.Reset()
	fn(current, w)
}

func (w *Writer) Value(v any, mode core.WritingMode) {
	expr.Parse(v).Render(w, w.dial, mode)
}

// --- Internal writer

func (w *Writer) New() core.InternalWriter {
	return &Writer{dial: w.dial}
}

func (w *Writer) ToString() string {
	return w.sb.String()
}

func (w *Writer) Reset() {
	w.sb.Reset()
}
