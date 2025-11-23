package core

type Config struct {
	Dialect   string
	ColWrite  ColumnWriter
	ColumnTag string
	Limit     int
	MaxLimit  int
}

type ColumnWriter struct {
	Select  ColumnType
	Join    ColumnType
	Where   ColumnType
	GroupBy ColumnType
	Having  ColumnType
	OrderBy ColumnType
}
