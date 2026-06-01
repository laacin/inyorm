package api

// --- ColDecl

type ColDecl interface {
	String(name string) ColDeclNext
	Int(name string) ColDeclNext
	Float(name string) ColDeclNext
	Bool(name string) ColDeclNext
}

type ColDeclNext interface {
	PrimaryKey() ColDeclNext
	AutoIncrement() ColDeclNext
	Unique() ColDeclNext
	Nullable() ColDeclNext
	Default(v any) ColDeclNext
}

// --- ConsDecl

type ConsDecl interface {
	PrimaryKey(on ...string)
	ForeignKey(on string) ForeignKey
	Check(ident any) Cond
}

type ForeignKey interface {
	To(col, table string) ForeignKeyNext
}

type ForeignKeyNext interface {
	OnDel(key string) ForeignKeyNext
	OnUpd(key string) ForeignKeyNext
}
