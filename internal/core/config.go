package core

type Config struct {
	Dialect   string
	AutoPh    AutoPlaceholder
	ColWrite  ColumnWriter
	ColumnTag string
	Limit     int
	MaxLimit  int
}

type AutoPlaceholder struct {
	Insert bool
	Update bool
	Where  bool
	Having bool
	Join   bool
}

type ColumnWriter struct {
	Select  ColumnType
	Join    ColumnType
	Where   ColumnType
	GroupBy ColumnType
	Having  ColumnType
	OrderBy ColumnType
}
