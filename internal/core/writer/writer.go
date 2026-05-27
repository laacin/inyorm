package writer

import (
	"strings"

	"github.com/laacin/inyorm/internal/core"
)

type Writer struct {
	sb     strings.Builder
	parser core.ValueParser
}

func New(parser core.ValueParser) *Writer {
	return &Writer{parser: parser}
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
	w.parser.Render(w, v, mode)
}

// --- Internal writer

func (w *Writer) New() core.InternalWriter {
	return &Writer{parser: w.parser}
}

func (w *Writer) ToString() string {
	return w.sb.String()
}

func (w *Writer) Reset() {
	w.sb.Reset()
}
