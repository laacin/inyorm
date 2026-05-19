package api

type TableBuilder interface {
	Text(name string) ColDecl
	Int(name string) ColDecl
	Float(name string) ColDecl
	Bool(name string) ColDecl

	Cons() ConsDecl
}

type ColDecl interface {
	PrimaryKey() ColDecl
	AutoIncrement() ColDecl
	Unique() ColDecl
	Nullable() ColDecl
}

type ConsDecl interface {
	Index(col string)
	ForeignKey(col string) ForeignKey
	Check(ident any) Condition
	Default(col string) Default
}

// --- Dependencies
type ForeignKey interface {
	To(col, table string) ForeignKeyNext
}

type ForeignKeyNext interface {
	OnDel(key string) ForeignKeyNext
	OnUpd(key string) ForeignKeyNext
}

type Default interface{ Value(value any) }
type Check interface{ On(ident any) Condition }
