package entity

type Table struct{ Value string }

func (*Table) Kind() ValueKind { return ValueTable }

func (t *Table) Write(w Writer, dial ValueWriter, mode WritingMode) {
	dial.WriteTable(w, t)
}
