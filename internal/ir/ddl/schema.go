package ddl

import "strings"

type ColumnMeta struct {
	Name          string
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Index         bool
	NotNull       bool
	Foreign       ForeignKey
}

type ForeignKey struct {
	Declared bool
	Table    string
	Column   string
	OnDelete OnAction
	OnUpdate OnAction
}

type OnAction int

const (
	OnActionUnset OnAction = iota
	OnActionCascade
	OnActionSetNull
	OnActionDefault
	OnActionRestrict
	OnActionNoAction
)

func ParseTag(tag string) ColumnMeta {
	result := ColumnMeta{}

	for seg := range strings.SplitSeq(tag, KeySep) {
		handleSegment(seg, &result)
	}

	return result
}

func handleSegment(seg string, ptr *ColumnMeta) {
	seg = strings.TrimSpace(seg)
	if seg == "" {
		return
	}

	splits := strings.SplitN(seg, ValSep, 2)
	key := strings.ToLower(strings.TrimSpace(splits[0]))

	if has(pkSet, key) {
		ptr.PrimaryKey = true
		return
	}

	if has(aiSet, key) {
		ptr.AutoIncrement = true
		return
	}

	if has(uqSet, key) {
		ptr.Unique = true
		return
	}

	if has(nnSet, key) {
		ptr.NotNull = true
		return
	}

	if has(idxSet, key) {
		ptr.Index = true
		return
	}

	if len(splits) < 2 {
		return
	}

	value := strings.TrimSpace(splits[1])

	if has(cSet, key) {
		ptr.Name = value
		return
	}

	if has(fkSet, key) {
		refs := strings.SplitN(value, ValSep, 2)
		if len(refs) < 2 {
			return
		}

		ptr.Foreign.Declared = true
		ptr.Foreign.Table = strings.TrimSpace(refs[0])
		ptr.Foreign.Column = strings.TrimSpace(refs[1])
		return
	}

	if has(odSet, key) {
		ptr.Foreign.OnDelete = setAct(strings.ToLower(value))
		return
	}

	if has(ouSet, key) {
		ptr.Foreign.OnUpdate = setAct(strings.ToLower(value))
		return
	}
}

func has(set map[string]struct{}, key string) bool {
	_, ok := set[key]
	return ok
}

func setAct(v string) OnAction {
	switch v {
	case ActCascade:
		return OnActionCascade
	case ActSetNull:
		return OnActionSetNull
	case ActDefault:
		return OnActionDefault
	case ActRestrict:
		return OnActionRestrict
	case ActNoAction:
		return OnActionNoAction
	default:
		return OnActionUnset
	}
}
