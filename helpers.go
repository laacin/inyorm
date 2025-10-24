package inyorm

type TableFor[T any] struct {
	Name        string
	PrimaryKey  string
	ForeignKeys map[string]string
	Field       T
}

func (t *TableFor[T]) TableName() string        { return t.Name }
func (t *TableFor[T]) PKey() string             { return t.PrimaryKey }
func (t *TableFor[T]) FKeys() map[string]string { return t.ForeignKeys }
