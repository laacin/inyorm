package ddl

import (
	"github.com/laacin/inyorm/internal/api"
	"github.com/laacin/inyorm/internal/core"
)

type ForeignKey struct {
	Col      string
	ToTable  string
	ToCol    string
	OnDelete OnAction
	OnUpdate OnAction
}

// start

func (b *ForeignKey) Start(col string) *ForeignKey {
	b.Col = col
	return b
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

func (b *ForeignKey) Build(w core.InternalWriter, dial TableWriter) {
	dial.WriteConsForeignKey(w, b)
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
