package core

type Builder struct {
	Mapper Mapper
	Param  ParamStore
	Attach AttachmentConfig
}

// --- Attachments

type AttachmentConfig struct {
	MainRef    string
	UseAliases bool
}

// --- Mapper

// tag rules
const (
	TAG = "inyorm"

	TagKeySep = ","
	TagValSep = ":"

	TagKeySkip = "skip"
	TagKeyCol  = "col"
)

type Mapper interface {
	ReadCols(...any) []string
	ReadValues(cols []string, v any) ([]any, error)

	Scan(rows Rows, bind any) error
}

type ParamStore interface {
	Store(v any)

	Lazy(id string)
	LazyObj(cols []string)

	Fill(id string, v any)
	FillObj(func(cols []string) []any) // Objs loads must be in orders

	LastIndex(idx int) ParamIndex
	Values() ([]any, error)
}

type ParamIndex struct {
	ID  string
	Num int
}
