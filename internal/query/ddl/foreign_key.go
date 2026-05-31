package ddl

import "github.com/laacin/inyorm/internal/api"

type ForeignKey struct {
	Col      string
	ToTable  string
	ToCol    string
	OnDelete OnAction
	OnUpdate OnAction
}

func NewForeignKey(col string) *ForeignKey {
	return &ForeignKey{Col: col}
}

// --- PUB API

func (b *ForeignKey) To(col, table string) api.ForeignKeyNext {
	b.ToTable = table
	b.ToCol = col
	return b
}

func (b *ForeignKey) OnDel(key string) api.ForeignKeyNext {
	b.OnDelete = setOnAct(key)
	return b
}
func (b *ForeignKey) OnUpd(key string) api.ForeignKeyNext {
	b.OnUpdate = setOnAct(key)
	return b
}

// --- Build

func (b *ForeignKey) Build() error {
	return nil
}

// --- Tools

type OnAction int

const (
	OnActionUnset OnAction = iota
	OnActionCascade
	OnActionSetNull
	OnActionDefault
	OnActionRestrict
	OnActionNoAction
)

func setOnAct(key string) OnAction {
	switch key {
	case "cascade":
		return OnActionCascade
	case "setnull":
		return OnActionSetNull
	case "default":
		return OnActionDefault
	case "restrict":
		return OnActionRestrict
	case "noaction":
		return OnActionNoAction
	default:
		return OnActionUnset
	}
}
