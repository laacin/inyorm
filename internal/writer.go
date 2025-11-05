package internal

import "strings"

type Writer struct {
	sb      strings.Builder
	ph      *Placeholder
	aliases *Alias
}

func (w *Writer) Write(v string) { w.sb.WriteString(v) }
func (w *Writer) Char(v byte)    { w.sb.WriteByte(v) }

func (w *Writer) Value(v any, placeholder bool) {
	if placeholder {
		w.sb.WriteString(w.ph.next(v))
		return
	}
	w.sb.WriteString(Sqlify(v))
}

func (w *Writer) InferColumn(v any, def bool) {
	switch val := v.(type) {
	case Column:
		w.Column(val, def)
	default:
		w.sb.WriteString(Sqlify(v))
	}
}

func (w *Writer) Column(v Column, def bool) {
	switch v.Typ {
	case NormalCol:
		alias := w.aliases.Get(v.Table)
		w.sb.WriteByte(alias)
		w.sb.WriteByte('.')
		w.sb.WriteString(v.Value)
	case CustomCol:
		if def {
			w.sb.WriteString(v.Value)
			w.sb.WriteString(" AS ")
		}
		w.sb.WriteString(v.Alias)
	case KeywordCol:
		w.sb.WriteString(v.Value)
	}
}

func (w *Writer) Table(v string) {
	w.sb.WriteString(v)
	w.sb.WriteByte(' ')
	w.sb.WriteByte(w.aliases.Get(v))
}

func (w *Writer) Read() string { return w.sb.String() }
func (w *Writer) Reset()       { w.sb.Reset() }
