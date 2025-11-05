package internal

type ColumnType int

const (
	NormalCol ColumnType = iota
	CustomCol
	KeywordCol
)

type Column struct {
	Typ   ColumnType
	Value string
	Alias string
	Table string
}
