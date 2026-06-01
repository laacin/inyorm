package api

// --- Select

type Select interface {
	Select(vals ...any) SelectNext
}

type SelectNext interface {
	Distinct()
}

// --- From

type From interface {
	From(v any)
}

type Into interface {
	Into(v any)
}

// --- Join

type Join interface {
	Join(v any) JoinNext
}

type JoinNext interface {
	On(ident any) Cond
	Left() JoinEnd
	Full() JoinEnd
	Cross()
}

type JoinEnd interface {
	On(ident any) Cond
}

// --- Where

type Where interface {
	Where(ident any) Cond
}

// --- Group By

type GroupBy interface {
	GroupBy(vals ...any)
}

// --- Having

type Having interface {
	Having(ident any) Cond
}

// --- Order By

type OrderBy interface {
	OrderBy(v any) OrderByNext
}

type OrderByNext interface {
	Desc()
}

// --- Limit

type Limit interface {
	Limit(v int)
}

// --- Offset

type Offset interface {
	Offset(v int)
}

// --- Insert

type Insert interface {
	Insert(ref ...any) Ignore
}

type OnConflict interface {
	OnConflict(ident ...any) OnConflictNext
}

type OnConflictNext interface {
	DoNothing()
	DoUpdate(ident ...any)
}

// --- Update

type Update interface {
	Update(ref ...any) Ignore
}

// --- Delete

type Delete interface {
	Delete()
}

// --- EXTRA

type Ignore interface {
	Ignore(ignore ...any)
}

type Values interface {
	Values(v any)
}
