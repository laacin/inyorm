package api

type TableBuilder interface {
	Text(name string) ColDecl
	Int(name string) ColDecl
	Float(name string) ColDecl
	Bool(name string) ColDecl

	ForeignKey(col string) ForeignKey
	Check(ident any) Condition
}

type ColDecl interface {
	PrimaryKey() ColDecl
	AutoIncrement() ColDecl
	Unique() ColDecl
	Nullable() ColDecl
	Default(value any) ColDecl
}

// --- Dependencies
type ForeignKey interface {
	To(col, table string) ForeignKeyNext
}

type ForeignKeyNext interface {
	OnDel(key string) ForeignKeyNext
	OnUpd(key string) ForeignKeyNext
}
