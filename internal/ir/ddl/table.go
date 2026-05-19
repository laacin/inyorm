package ddl

type TableDecl struct {
	Name string
	Cols []ColDecl
	Cons []ConsDecl[any]
}
