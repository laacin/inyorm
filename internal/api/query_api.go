package api

type SelectQuery interface {
	Select
	From
	Join
	Where
	GroupBy
	Having
	OrderBy
	Limit
	Offset
}

type InsertQuery interface {
	Insert
	Into
	Values
	OnConflict
}

type UpdateQuery interface {
	Update
	Into
	Values
	Where
}

type DeleteQuery interface {
	Delete
	From
	Where
}

type CreateTableQuery interface {
	TableName(name string)

	ColDecl
	ConsDecl
}
