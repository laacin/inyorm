package ddl

import "github.com/laacin/inyorm/internal/api"

type ConsKind int

const (
	ConsKindIndex ConsKind = iota
	ConsKindForeignKey
	ConsKindCheck
)

type ConsDecl[T any] struct {
	Kind   ConsKind
	Table  string
	Column string
	Value  T
}

func (c *ConsDecl[T]) IsIndex() (*ConsDecl[ConsIndex], bool) {
	if c.Kind == ConsKindIndex {
		return assertConsDecl[ConsIndex](c), true
	}
	return nil, false
}

func (c *ConsDecl[T]) IsForeignKey() (*ConsDecl[ConsForeignKey], bool) {
	if c.Kind == ConsKindForeignKey {
		return assertConsDecl[ConsForeignKey](c), true
	}
	return nil, false
}

func (c *ConsDecl[T]) IsCheck() (*ConsDecl[ConsCheck], bool) {
	if c.Kind == ConsKindCheck {
		return assertConsDecl[ConsCheck](c), true
	}
	return nil, false
}

func assertConsDecl[T, K any](c *ConsDecl[K]) *ConsDecl[T] {
	return &ConsDecl[T]{
		Kind:   c.Kind,
		Table:  c.Table,
		Column: c.Column,
		Value:  *any(c.Value).(*T),
	}
}

// --- Types

type ConsIndex struct{}

type ConsForeignKey struct {
	ToTable  string
	ToColumn string
	OnDelete OnAction
	OnUpdate OnAction
}

type ConsDefault struct{ Value any }

type ConsCheck struct{ Cond api.Condition }

// --- Dependencies
type OnAction int

const (
	OnActionUnset OnAction = iota
	OnActionCascade
	OnActionSetNull
	OnActionDefault
	OnActionRestrict
	OnActionNoAction
)

func SetOnAct(key string) OnAction {
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
